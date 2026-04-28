---
icon: lucide/plug
---

# API 参考

Agent-Sandbox 提供两类 API 接口：**原生 REST API** 和 **E2B 兼容 API**。所有接口都需要在 `Authorization` 头中传递 API 密钥。

## 认证

所有 API 请求必须包含 API 密钥：

```bash
Authorization: Bearer <your-api-key>
```

默认系统令牌：`sys-2492a85b10ed4cb083b2c76b181eac96`

---

## 原生 REST API

基础路径：`/api/v1`

### 创建沙箱

```
POST /api/v1/sandbox
```

**请求体：**

| 字段 | 类型 | 必填 | 描述 |
|-------|------|----------|-------------|
| name | string | 否 | 沙箱名称。省略时自动生成。小写字母、数字和 `-`，最长 50 字符。 |
| image | string | 否 | 容器镜像覆盖。 |

**示例：**

```bash
curl -X POST http://<host>/api/v1/sandbox \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <api-key>" \
  -d '{"name":"sandbox-01"}'
```

**响应 (200)：**

```json
{
  "code": "0",
  "data": "Sandbox sandbox-01 created successfully"
}
```

### 列出沙箱

```
GET /api/v1/sandbox
```

**示例：**

```bash
curl http://<host>/api/v1/sandbox \
  -H "Authorization: Bearer <api-key>"
```

**响应 (200)：** 沙箱对象数组。

### 获取沙箱

```
GET /api/v1/sandbox/{name}
```

**路径参数：**

| 参数 | 类型 | 描述 |
|-----------|------|-------------|
| name | string | 沙箱名称 |

**响应 (200)：** 沙箱详情对象。

### 删除沙箱

```
DELETE /api/v1/sandbox/{name}
```

**示例：**

```bash
curl -X DELETE http://<host>/api/v1/sandbox/sandbox-01 \
  -H "Authorization: Bearer <api-key>"
```

**响应 (200)：**

```json
{
  "code": "0",
  "data": "Sandbox sandbox-01 deleted successfully"
}
```

## E2B 兼容 API

基础路径：`/e2b/v1`

这些接口兼容 E2B 协议和 SDK。

### 创建沙箱

```
POST /e2b/v1/sandboxes
```

**请求体：**

| 字段 | 类型 | 必填 | 描述 |
|-------|------|----------|-------------|
| templateID | string | 是 | 模版标识符（如 `code-interpreter`）。`-v1` 等版本后缀会自动去除。仅允许字母数字、连字符和点。 |
| timeout | int | 否 | 存活时间（秒）。 |
| envVars | object | 否 | 环境变量键值对。 |
| metadata | object | 否 | 自定义元数据键值对。 |
| secure | bool | 否 | 启用仅 HTTPS 模式。 |
| allowInternetAccess | bool | 否 | 允许出站网络访问（默认 true）。 |
| autoPause | bool | 否 | 超时后自动暂停而非终止。 |
| network | object | 否 | 网络配置，包含 `allowOut`、`denyOut`、`allowPublicTraffic`、`maskRequestHost`。 |
| mcp | object | 否 | MCP 配置。 |

**示例：**

```bash
curl -X POST http://<host>/e2b/v1/sandboxes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <api-key>" \
  -d '{"templateID":"code-interpreter","timeout":300}'
```

**响应 (201)：**

```json
{
  "sandboxID": "e94466d4e94466d4",
  "clientID": "client-id-x",
  "envdVersion": "0.1.1",
  "envdAccessToken": "e94466d4e94466d4",
  "templateID": "code-interpreter",
  "trafficAccessToken": "e94466d4e94466d4",
  "domain": "example.domain.com",
  "metadata": { "name": "sandbox-abc123" },
  "cpuCount": 2000,
  "memoryMB": 4096,
  "diskSizeMB": 5120,
  "startedAt": "2026-01-27T10:00:00Z",
  "endAt": "2026-01-27T10:05:00Z",
  "state": "running"
}
```

### 获取沙箱

```
GET /e2b/v1/sandboxes/{sandboxID}
```

**路径参数：**

| 参数 | 类型 | 描述 |
|-----------|------|-------------|
| sandboxID | string | 沙箱标识符 |

**响应 (200)：** 与创建沙箱响应相同。

**错误 (404)：**

```json
{
  "code": 404,
  "message": "sandbox <sandboxID> not found"
}
```

### 列出沙箱

```
GET /e2b/v1/v2/sandboxes
```

**响应 (200)：** 沙箱对象数组。

### 连接沙箱

```
POST /e2b/v1/sandboxes/{sandboxID}/connect
```

**请求体：**

| 字段 | 类型 | 必填 | 描述 |
|-------|------|----------|-------------|
| timeout | int32 | 否 | 从现在起沙箱过期的超时秒数。 |

**响应 (201)：** 沙箱对象（与获取沙箱响应格式相同）。

### 删除沙箱

```
DELETE /e2b/v1/sandboxes/{sandboxID}
```

**响应 (200)：** 成功时返回空。

---

## 错误格式

所有 API 错误遵循此 JSON 结构：

```json
{
  "code": 400,
  "message": "人类可读的错误描述"
}
```

常见 HTTP 状态码：

| 状态码 | 含义 |
|--------|---------|
| 400 | 错误请求（无效输入、缺少参数） |
| 401 | 未授权（缺少或无效 API 密钥） |
| 404 | 沙箱未找到 |
| 500 | 内部服务器错误 |

---


## 沙箱工具 REST API

### 终端

```
POST /api/v1/terminal/sandbox/{name}
```

在沙箱中执行命令。

### 文件

```
GET    /api/v1/sandbox/files/{name}                # 列出文件
POST   /api/v1/sandbox/files/{name}/upload          # 上传文件
DELETE /api/v1/sandbox/files/{name}                  # 删除文件
GET    /api/v1/sandbox/files/{name}/download         # 下载文件
```

### 日志

```
GET /api/v1/logs/sandbox/{name}            # 获取沙箱日志
```

### 列出事件

```
GET /api/v1/events
```

**响应 (200)：** 沙箱事件数组。

## 配置 API

```
GET  /api/v1/config/templates              # 获取模版配置
POST /api/v1/config/templates              # 保存模版配置
GET  /api/v1/config/sandbox-template       # 获取沙箱模版配置
POST /api/v1/config/sandbox-template       # 保存沙箱模版配置
```

## 池 API

```
GET    /api/v1/pool                         # 列出池
GET    /api/v1/pool/sandbox/{name}          # 列出池沙箱
DELETE /api/v1/pool/{name}                  # 删除池
```


---

## 其他端点

| 路径 | 描述 |
|------|-------------|
| `/mcp` | MCP 服务端点（Streamable HTTP） |
| `/sandbox/{name}/` | 沙箱容器路由 |
| `/sandboxes/router/{sandboxID}/{port}/` | E2B 风格沙箱端口路由 |
| `/healthz` | 健康检查，返回 `{"status":"ok","version":"..."}` |
| `/ui/` | Web UI 静态文件 |