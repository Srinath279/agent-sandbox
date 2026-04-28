---
icon: lucide/settings
---

# 环境变量

Agent-Sandbox 通过 [envconfig](https://github.com/kelseyhightower/envconfig) 使用环境变量进行配置。本文档记录所有可用的配置选项。

---

## 服务配置

### `SERVER_ADDR`

服务监听地址。

- **默认值：** `0.0.0.0:10000`
- **示例：** `0.0.0.0:8080`

### `API_BASE_URL`

计算为 `/api/{API_VERSION}`。通常不需要覆盖。

- **默认值：** `/api/v1`

---

## 认证

### `SYSTEM_TOKEN`

系统管理员令牌，拥有所有资源的完全访问权限。

- **默认值：** `sys-2492a85b10ed4cb083b2c76b181eac96`
- **示例：** `sys-your-custom-admin-token`

### `API_TOKENS_RAW`

逗号分隔的额外 API 令牌列表。`SYSTEM_TOKEN` 会自动添加到列表前面。

- **默认值：** `testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f`
- **示例：** `user1-abc123,user2-def456,user3-ghi789`
- **注意：** 令牌长度至少 5 个字符。

令牌在启动时解析和验证。无效或空令牌会被忽略。

---

## Kubernetes 配置

### `SANDBOX_NAMESPACE`

创建沙箱 Pod 和 ReplicaSet 的 Kubernetes 命名空间。

- **默认值：** `default`
- **示例：** `agent-sandbox`

### `CONFIGMAP_NAME`

用于存储模版和沙箱模版配置的 ConfigMap 名称。

- **默认值：** `agent-sandbox`
- **示例：** `my-sandbox-config`

### `ENV_NAME`

用于部署标识的环境名称标签。

- **默认值：** `dev`
- **示例：** `production`、`staging`

---

## 沙箱默认值

### `SANDBOX_DEFAULT_IMAGE`

创建沙箱时未指定模版使用的默认容器镜像。

- **默认值：** `ghcr.io/agent-infra/sandbox:latest`
- **示例：** `myregistry/my-sandbox:v1.0.0`

### `SANDBOX_DEFAULT_TEMPLATE`

未指定时使用的默认模版名称。

- **默认值：** `aio`
- **示例：** `code-interpreter`


---

## 配置示例

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
```

---

## 参考

| 变量 | 默认值 | 描述 |
|----------|---------|-------------|
| `SERVER_ADDR` | `0.0.0.0:10000` | 服务监听地址 |
| `SYSTEM_TOKEN` | `sys-2492a85b10ed4cb083b2c76b181eac96` | 系统管理员令牌 |
| `API_TOKENS_RAW` | `testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f` | 逗号分隔的用户令牌 |
| `SANDBOX_NAMESPACE` | `default` | 沙箱的 Kubernetes 命名空间 |
| `CONFIGMAP_NAME` | `agent-sandbox` | 配置存储的 ConfigMap 名称 |
| `SANDBOX_DEFAULT_IMAGE` | `ghcr.io/agent-infra/sandbox:latest` | 默认沙箱镜像 |
| `SANDBOX_DEFAULT_TEMPLATE` | `aio` | 默认模版名称 |
| `ENV_NAME` | `dev` | 环境名称标签 |