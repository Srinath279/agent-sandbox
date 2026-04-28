---
icon: lucide/play
---

# 快速开始

本指南将带你完成 Agent-Sandbox 的部署，并通过 E2B SDK 和 REST API 创建你的第一个沙箱。

## 前置条件

- Kubernetes 集群（版本 1.26 或更高）
- 已配置 `kubectl` 可访问你的集群
- （可选）`e2b-code-interpreter` Python SDK：`pip install e2b-code-interpreter`

## 1. 部署 Agent-Sandbox

应用安装清单：

```bash
kubectl create namespace agent-sandbox
kubectl apply -n agent-sandbox -f install.yaml
```

!!! Tip

    下载 `install.yaml`：https://github.com/agent-sandbox/agent-sandbox/blob/main/install.yaml


[install.yaml](https://github.com/agent-sandbox/agent-sandbox/blob/main/install.yaml) 包含完整的部署配置和 `Agent-Sandbox` 控制器组件。

## 2. 暴露服务

创建 Ingress：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: agent-sandbox
  namespace: agent-sandbox
spec:
  ingressClassName: ingress-nginx
  rules:
  - host: agent-sandbox.your-host.com
    http:
      paths:
      - backend:
          service:
            name: agent-sandbox
            port:
              number: 80
        path: /
```

或本地测试使用端口转发：

```bash
kubectl port-forward -n agent-sandbox svc/agent-sandbox 8080:80
```

## 3. 验证部署

```bash
curl https://agent-sandbox.your-host.com/healthz
```

预期响应：

```json
{"status":"ok","version":"xxx"}
```

## 4. 认证

获取你的 API 密钥。默认系统令牌是：

```
sys-2492a85b10ed4cb083b2c76b181eac96
```

所有 API 请求需要 `Authorization` 头：

```bash
Authorization: Bearer <your-api-key>
```

---

## 通过 E2B SDK 创建沙箱

E2B SDK 是 Python 应用的推荐方式。

### 安装 SDK

```bash
pip install e2b-code-interpreter
```

### 配置环境

```python
import os

# 将 SDK 指向你的 Agent-Sandbox 实例
os.environ['E2B_API_URL'] = 'https://agent-sandbox.your-host.com/e2b/v1'
os.environ['E2B_API_KEY'] = 'sys-2492a85b10ed4cb083b2c76b181eac96'
os.environ['E2B_DOMAIN'] = 'agent-sandbox.your-host.com'
```

#### 本地开发配置
!!! Warning

    对于本地或开发环境（无 HTTPS），SDK 的默认 API URL 无法工作。使用以下配置连接本地实例：

```python
def local():
    os.environ['E2B_DEBUG'] = "true"
    os.environ['E2B_API_URL'] = 'http://localhost:10000/e2b/v1'
    os.environ['E2B_API_KEY'] = 'testuser-aef134ef-7aa1-945e-9399-7df9a4ad0c3f'
    os.environ['E2B_DOMAIN'] = 'localhost:10000'

    def __connection_config_get_host(_, sandbox_id: str, sandbox_domain: str, port: int) -> str:
        return f"{sandbox_domain}/sandboxes/router/{sandbox_id}/{port}"
    from e2b import ConnectionConfig
    ConnectionConfig.get_host = __connection_config_get_host
```



### 创建和使用沙箱

```python
from e2b_code_interpreter import Sandbox

# 创建沙箱
sandbox = Sandbox.create(template='code-interpreter', timeout=300)

print(f"Sandbox ID: {sandbox.sandbox_id}")

# 执行代码
execution = sandbox.run_code("print('Hello from Agent-Sandbox!')")
print(execution.logs)

# 运行 shell 命令
result = sandbox.commands.run("ls -la")
print(result.stdout)

# 清理
sandbox.kill()
```

### 连接到已有沙箱

```python
from e2b_code_interpreter import Sandbox

# 通过 ID 重连
sandbox = Sandbox.connect("existing-sandbox-id")
print(f"Connected to: {sandbox.sandbox_id}")
```

---

## 通过 REST API 创建沙箱

对于非 Python 环境，直接使用 REST API。

### 原生 REST API

```bash
curl -X POST http://localhost:8080/api/v1/sandbox \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sys-2492a85b10ed4cb083b2c76b181eac96" \
  -d '{"name":"my-sandbox"}'
```

响应：

```json
{
  "code": "0",
  "data": "Sandbox my-sandbox created successfully"
}
```

### E2B 兼容 API

```bash
curl -X POST http://localhost:8080/e2b/v1/sandboxes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sys-2492a85b10ed4cb083b2c76b181eac96" \
  -d '{"templateID":"code-interpreter","timeout":300}'
```

响应：

```json
{
  "sandboxID": "e94466d4e94466d4",
  "templateID": "code-interpreter",
  "clientID": "client-id-x",
  "envdVersion": "0.1.1",
  "envdAccessToken": "e94466d4e94466d4",
  "trafficAccessToken": "e94466d4e94466d4",
  "domain": "localhost:8080",
  "metadata": {"name": "sandbox-abc123"},
  "cpuCount": 2000,
  "memoryMB": 4096,
  "diskSizeMB": 5120,
  "startedAt": "2026-01-27T10:00:00Z",
  "endAt": "2026-01-27T10:05:00Z",
  "state": "running"
}
```

### 列出沙箱

```bash
curl http://localhost:8080/api/v1/sandbox \
  -H "Authorization: Bearer sys-2492a85b10ed4cb083b2c76b181eac96"
```

### 删除沙箱

```bash
curl -X DELETE http://localhost:8080/api/v1/sandbox/my-sandbox \
  -H "Authorization: Bearer sys-2492a85b10ed4cb083b2c76b181eac96"
```

---