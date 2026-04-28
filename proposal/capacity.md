# Sandbox Creation Rate Limit Design

## 1. Requirements

Add rate limiting for sandbox creation APIs:

- Limit concurrent sandbox creation requests per user.
- Return `429 Too Many Requests` immediately when concurrency is exceeded. Do not queue or wait.
- Limit total sandbox capacity per user.
- Return `413 Payload Too Large` when capacity is exceeded, including current and max counts.
- Rate limit failures must not block normal sandbox creation. Unexpected limiter or capacity-check errors degrade gracefully and allow the request.
- Keep the existing request and response structures unchanged.
- Keep the rate limit module independent.

The rate limit checks apply to both creation paths:

- Native API: `POST /api/v1/sandbox`
- E2B API sandbox creation path

## 2. Core Concepts

| Concept | Meaning |
| --- | --- |
| Concurrency | Number of in-flight sandbox creation requests for the same user. |
| Capacity | Number of existing non-pool sandboxes owned by the same user. |

## 3. Limit Behavior

| Limit | Check time | Limit target | Exceeded behavior | HTTP status |
| --- | --- | --- | --- | --- |
| Concurrency | Before sandbox creation starts | In-flight creation requests | Reject immediately | 429 |
| Capacity | Before concurrency slot acquisition | Existing sandbox count | Reject immediately | 413 |

Capacity uses `413` instead of `403` to avoid conflicting with authentication or authorization failures.

## 4. Request Flow

```text
Create request
  -> Check user capacity
      -> exceeded: return 413
      -> check failed: log warning and continue
  -> Try acquire concurrency slot
      -> exceeded: return 429
      -> check failed: log warning and continue
  -> Create sandbox
  -> Release concurrency slot if acquired
```

## 5. Configuration

Rate limit config is stored directly in `Config`.

```go
type RateLimitConfig struct {
    Enabled        bool `split_words:"true" default:"true"`
    MaxConcurrency int  `split_words:"true" default:"3"`
    MaxSandbox     int  `split_words:"true" default:"100"`
}

type UserRateLimitConfig struct {
    User           string `json:"user"`
    MaxConcurrency int    `json:"max_concurrency"`
    MaxSandbox     int    `json:"max_sandbox"`
}

type Config struct {
    RateLimit         RateLimitConfig       `split_words:"true"`
    RateLimitUsersRaw string                `split_words:"true" required:"false"`
    RateLimitUsers    []UserRateLimitConfig `json:"rateLimitUsers"`
}
```

### 5.1 Environment variables

The default rate limit config is loaded by `envconfig` from environment variables:

| Environment variable | Config field | Default | Meaning |
| --- | --- | --- | --- |
| `RATE_LIMIT_ENABLED` | `Config.RateLimit.Enabled` | `true` | Enables rate limiting. |
| `RATE_LIMIT_MAX_CONCURRENCY` | `Config.RateLimit.MaxConcurrency` | `3` | Default max concurrent creation requests per user. `0` means unlimited. |
| `RATE_LIMIT_MAX_SANDBOX` | `Config.RateLimit.MaxSandbox` | `100` | Default max existing non-pool sandboxes per user. `0` means unlimited. |
| `RATE_LIMIT_USERS` | `Config.RateLimitUsersRaw` | empty | JSON array for user-specific quotas. Invalid JSON is ignored with a warning. |

Example:

```bash
RATE_LIMIT_ENABLED=true
RATE_LIMIT_MAX_CONCURRENCY=5
RATE_LIMIT_MAX_SANDBOX=200
RATE_LIMIT_USERS='[{"user":"testuser-aaa","max_concurrency":10,"max_sandbox":500},{"user":"vip-user-xxx","max_concurrency":20,"max_sandbox":1000}]'
```

`RATE_LIMIT_USERS` must be a JSON array. Each entry uses the same shape as `UserRateLimitConfig`.

Runtime user quota config can also be updated through the application ConfigMap key `config-runtime`. The service only watches and reads this key; it does not create or overwrite it. If the key is missing, environment defaults are used. Changes take effect without restarting the server.

Example ConfigMap value:

```json
{
  "rate_limit_users_raw": "[{\"user\":\"testuser-aaa\",\"max_concurrency\":10,\"max_sandbox\":500}]"
}
```

### 5.2 User-specific config

User-specific config overrides defaults. A zero value in a user-specific config means use the default.

Example:

```json
[
  {
    "user": "testuser-aaa",
    "max_concurrency": 10,
    "max_sandbox": 500
  },
  {
    "user": "vip-user-xxx",
    "max_concurrency": 20,
    "max_sandbox": 1000
  }
]
```

Config validation rules:

- Negative default values are reset to the built-in defaults.
- Empty user entries are ignored.
- Negative user-specific values are reset to `0`, which means use the default.

## 6. Rate Limiter Module

The limiter lives in `pkg/ratelimit` and owns only in-memory concurrency state plus the sandbox controller reference. It reads `config.Cfg` at runtime so config changes can take effect without rebuilding the limiter.

```go
type sandboxController interface {
    CountByUser(user string) (int, error)
}

type UserLimiter struct {
    mu    sync.Mutex
    count int
}

type RateLimiter struct {
    defaultConfig *config.RateLimitConfig
    userConfigs   sync.Map
    userLimiters  sync.Map
    controller    sandboxController
}
```

Key methods:

| Method | Purpose |
| --- | --- |
| `CheckCapacity(user)` | Counts existing non-pool sandboxes for the user and reports capacity excess. |
| `TryAcquire(user)` | Attempts to acquire an in-memory concurrency slot without waiting. |
| `Release(user)` | Releases a previously acquired concurrency slot and never lets count go below zero. |
| `Stats(user)` | Returns a point-in-time status snapshot for status APIs. |
| `UserConfigs()` | Returns configured user quota snapshots for status APIs. |

## 7. Integration Points

### 7.1 Server initialization

`pkg/handler/handler.go` initializes the global limiter once during server setup:

```go
if config.Cfg.RateLimit.Enabled {
    ratelimit.Init(&config.Cfg.RateLimit, config.Cfg.RateLimitUsers, c)
    klog.Info("Rate limiter initialized")
} else {
    klog.Info("Rate limiter disabled")
}
```

### 7.2 Native sandbox creation

`pkg/handler/handlers.go` checks capacity first, then concurrency, before the existing create flow.

### 7.3 E2B sandbox creation

`pkg/api/e2b/sandbox.go` applies the same capacity and concurrency checks before creating the internal sandbox.

### 7.4 Capacity counting

`pkg/sandbox/controller.go` implements:

```go
func (s *Controller) CountByUser(user string) (int, error)
func (s *Controller) ListAllUsers() ([]string, error)
```

Capacity counting uses Kubernetes labels and only counts non-pool sandboxes:

- `owner=agent-sandbox`
- `sbx-user=<user>`
- `sbx-pool=false`

## 8. Error Responses

### 8.1 Concurrency exceeded

HTTP status: `429 Too Many Requests`

```json
{
  "code": "429",
  "error": "rate limit exceeded: too many concurrent sandbox creation requests (3/3), please retry later"
}
```

### 8.2 Capacity exceeded

HTTP status: `413 Payload Too Large`

```json
{
  "code": "413",
  "error": "capacity exceeded: you have reached the maximum sandbox limit (100/100), please delete some sandboxes before creating new ones"
}
```

Native handlers use `HTTPError` so the response body keeps the existing structure while the HTTP status can be `429` or `413`.

## 9. Status API

There is a single status API:

```text
GET /api/v1/ratelimit
```

It returns the full rate limit status for all users. There is no separate user-level status API and no `/ratelimit/all` endpoint.

The user list includes:

- Users with existing sandboxes.
- Users with explicit `RateLimitUsers` config even if they have no current sandbox.

Example response:

```json
{
  "code": "0",
  "data": {
    "default_config": {
      "max_concurrency": 3,
      "max_sandbox": 100
    },
    "users": [
      {
        "user": "testuser-aaa",
        "concurrency_active": 2,
        "concurrency_max": 10,
        "sandbox_current": 50,
        "sandbox_max": 500,
        "sandbox_usage_percent": 10
      },
      {
        "user": "vip-user-xxx",
        "concurrency_active": 0,
        "concurrency_max": 20,
        "sandbox_current": 0,
        "sandbox_max": 1000,
        "sandbox_usage_percent": 0
      }
    ]
  }
}
```

If the limiter is not initialized, the endpoint returns zero default config and an empty user list:

```json
{
  "code": "0",
  "data": {
    "default_config": {
      "max_concurrency": 0,
      "max_sandbox": 0
    },
    "users": []
  }
}
```

## 10. Graceful Degradation

| Scenario | Behavior |
| --- | --- |
| `GlobalLimiter == nil` | Skip create-time rate limit checks. |
| `Enabled == false` | Skip create-time rate limit checks. |
| `controller == nil` | Capacity check returns an internal error; create flow logs and continues. |
| `CountByUser()` returns error | Log warning and continue creation. |
| `TryAcquire()` returns error | Log warning and continue creation. |
| `Release()` called more than once | Do not panic and do not let count go below zero. |

## 11. Thread Safety

| Component | Protection |
| --- | --- |
| `UserLimiter.count` | `sync.Mutex` |
| `RateLimiter.userLimiters` | `sync.Map` |
| `RateLimiter.userConfigs` | `sync.Map` |
| `GlobalLimiter` | Initialized once at server startup and then read by request handlers. |

## 12. Test Plan

### 12.1 Unit tests

- Acquire and release concurrency slots.
- Return `acquired=false` when concurrency reaches the max.
- Return `exceeded=true` when capacity reaches the max.
- Treat `MaxConcurrency=0` as unlimited.
- Treat `MaxSandbox=0` as unlimited.
- Apply user-specific overrides.
- Use defaults when user-specific values are zero.
- Ensure repeated `Release()` does not panic or make counts negative.
- Validate negative and empty config handling.

### 12.2 Integration tests

- Concurrent sandbox creation returns `429` after the configured concurrency limit.
- Capacity limit returns `413` after the configured sandbox max.
- Capacity-check failures do not block creation.
- Concurrency-check failures do not block creation.
- Deleting sandboxes allows creation again after capacity drops below the max.
- `GET /api/v1/ratelimit` returns default config and all known/configured users.
