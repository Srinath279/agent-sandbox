---
icon: lucide/settings
---

# Environment Variables

Agent-Sandbox is configured via environment variables using [envconfig](https://github.com/kelseyhightower/envconfig). This page documents all available configuration options.

---

## Server Configuration

### `SERVER_ADDR`

Server listen address.

- **Default:** `0.0.0.0:10000`
- **Example:** `0.0.0.0:8080`

### `API_BASE_URL`

Computed as `/api/{API_VERSION}`. Not normally overridden.

- **Default:** `/api/v1`

---

## Authentication

### `SYSTEM_TOKEN`

System administrator token with full access to all resources.

- **Default:** `sys-2492a85b10ed4cb083b2c76b181eac96`
- **Example:** `sys-your-custom-admin-token`

### `API_TOKENS_RAW`

Comma-separated list of additional API tokens. The `SYSTEM_TOKEN` is automatically prepended.

- **Default:** `testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f`
- **Example:** `user1-abc123,user2-def456,user3-ghi789`
- **Note:** Tokens must be at least 5 characters long.

Tokens are parsed and validated at startup. Invalid or empty tokens are ignored.

---

## Kubernetes Configuration

### `SANDBOX_NAMESPACE`

Kubernetes namespace where sandbox pods and ReplicaSets are created.

- **Default:** `default`
- **Example:** `agent-sandbox`

### `CONFIGMAP_NAME`

Name of the ConfigMap used to store templates and sandbox template configuration.

- **Default:** `agent-sandbox`
- **Example:** `my-sandbox-config`

---

## Sandbox Defaults

### `SANDBOX_DEFAULT_IMAGE`

Default container image used when creating sandboxes without specifying a template.

- **Default:** `ghcr.io/agent-infra/sandbox:latest`
- **Example:** `myregistry/my-sandbox:v1.0.0`

### `SANDBOX_DEFAULT_TEMPLATE`

Default template name used when not specified.

- **Default:** `aio`
- **Example:** `code-interpreter`

---

## E2B Compatibility

These environment variables are used by E2B SDK clients to connect to Agent-Sandbox.

### `E2B_DOMAIN`

Domain used for sandbox URL generation.

- **Default:** `localhost`
- **Example:** `sandbox.example.com`
- **SDK Usage:** Sets the base domain for sandbox access URLs.

### `E2B_API_URL`

Base URL for E2B-compatible API endpoints.

- **Default:** (none, SDK uses default)
- **Example:** `http://localhost:10000/e2b/v1`
- **SDK Usage:** Override the API endpoint for E2B SDK calls.

### `E2B_API_KEY`

API key used by E2B SDK for authentication.

- **Default:** (none)
- **Example:** `sys-2492a85b10ed4cb083b2c76b181eac96`
- **SDK Usage:** Authentication token passed in `Authorization: Bearer` header.

---

## Sidecar: sucall-server

The `sucall-server` sidecar has its own environment variables.

### `SUCALL_ALLOWED_COMMAND_PREFIXES`

Comma-separated list of allowed command prefixes for the `sucall` wrapper.

- **Default:** `s3fs,fusermount`
- **Example:** `ls,s3fs,fusermount,mount,umount`
- **Note:** Commands are matched by exact prefix. `export` is always allowed as a built-in.

---

## Miscellaneous

### `ENV_NAME`

Environment name label for deployment identification.

- **Default:** `dev`
- **Example:** `production`, `staging`

---

## Example Configuration

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: agent-sandbox
  namespace: agent-sandbox
spec:
  template:
    spec:
      containers:
      - name: agent-sandbox
        env:
        - name: SYSTEM_TOKEN
          value: "sys-your-admin-token"
        - name: API_TOKENS_RAW
          value: "user1-abc123,user2-def456"
        - name: SANDBOX_NAMESPACE
          value: "agent-sandbox"
        - name: E2B_DOMAIN
          value: "sandbox.example.com"
```

---

## Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_ADDR` | `0.0.0.0:10000` | Server listen address |
| `SYSTEM_TOKEN` | `sys-2492a85b10ed4cb083b2c76b181eac96` | System admin token |
| `API_TOKENS_RAW` | `testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f` | Comma-separated user tokens |
| `SANDBOX_NAMESPACE` | `default` | Kubernetes namespace for sandboxes |
| `CONFIGMAP_NAME` | `agent-sandbox` | ConfigMap name for config storage |
| `SANDBOX_DEFAULT_IMAGE` | `ghcr.io/agent-infra/sandbox:latest` | Default sandbox image |
| `SANDBOX_DEFAULT_TEMPLATE` | `aio` | Default template name |
| `E2B_DOMAIN` | `localhost` | E2B SDK domain |
| `E2B_API_URL` | (none) | E2B API URL override |
| `E2B_API_KEY` | (none) | E2B SDK API key |
| `SUCALL_ALLOWED_COMMAND_PREFIXES` | `ls,s3fs,fusermount` | Allowed commands for sucall |
| `ENV_NAME` | `dev` | Environment name label |