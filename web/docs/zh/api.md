---
icon: lucide/plug
---

# API 参考

Agent-Sandbox 提供两类 API 接口：**原生 REST API** 和 **E2B 兼容 API**。所有接口均需在 `Authorization` 请求头携带 API Key。

## 鉴权

所有请求必须包含 API Key：

```bash
Authorization: Bearer <your-api-key>
```

默认系统用户 token：`sys-2492a85b10ed4cb083b2c76b181eac96`

---

## 原生 REST API

基础路径：`/api/v1`

### 创建沙箱

```
POST /api/v1/sandbox
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 沙箱名称，不填则自动生成。仅支持小写字母、数字和 `-`，最大 50 字符。 |
| image | string | 否 | 容器镜像覆盖配置。 |

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

| 参数 | 类型 | 说明 |
|------|------|------|
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

### 沙箱指标

```
POST /api/v1/sandbox/metrics
```

**响应 (200)：** 沙箱指标数据。

### 事件列表

```
GET /api/v1/events
```

**响应 (200)：** 沙箱事件数组。

### 终端

```
POST /api/v1/terminal/sandbox/{name}
```

在沙箱中执行命令。

```
GET /api/v1/terminal/sandbox/{name}/ws
```

WebSocket 终端流。

### 文件操作

```
GET    /api/v1/sandbox/files/{name}                # 列出文件
POST   /api/v1/sandbox/files/{name}/upload          # 上传文件
DELETE /api/v1/sandbox/files/{name}                  # 删除文件
GET    /api/v1/sandbox/files/{name}/download         # 下载文件
```

### 配置管理

```
GET  /api/v1/config/templates              # 获取模板配置
POST /api/v1/config/templates              # 保存模板配置
GET  /api/v1/config/sandbox-template       # 获取 sandbox-template 配置
POST /api/v1/config/sandbox-template       # 保存 sandbox-template 配置
```

### 池管理

```
GET    /api/v1/pool                         # 列出池
GET    /api/v1/pool/sandbox/{name}          # 列出池内沙箱
DELETE /api/v1/pool/{name}                  # 删除池
```

### 日志

```
GET /api/v1/logs/sandbox/{name}            # 获取沙箱日志
```

---

## E2B 兼容 API

基础路径：`/e2b/v1`

以下接口与 E2B 协议和 SDK 兼容。

### 创建沙箱

```
POST /e2b/v1/sandboxes
```

**请求体：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| templateID | string | 是 | 模板标识（如 `code-interpreter`）。版本后缀如 `-v1` 会自动去除。仅支持字母、数字、连字符和点。 |
| timeout | int | 否 | 存活时间（秒）。 |
| envVars | object | 否 | 环境变量键值对。 |
| metadata | object | 否 | 自定义元数据键值对。 |
| secure | bool | 否 | 启用 HTTPS-only 模式。 |
| allowInternetAccess | bool | 否 | 允许出站网络（默认 true）。 |
| autoPause | bool | 否 | 超时后自动暂停而非销毁。 |
| network | object | 否 | 网络配置，含 `allowOut`、`denyOut`、`allowPublicTraffic`、`maskRequestHost`。 |
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

| 参数 | 类型 | 说明 |
|------|------|------|
| sandboxID | string | 沙箱标识 |

**响应 (200)：** 与创建沙箱响应结构相同。

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

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| timeout | int32 | 否 | 从当前时间起算的超时秒数。 |

**响应 (201)：** 沙箱对象（结构与获取沙箱相同）。

### 删除沙箱

```
DELETE /e2b/v1/sandboxes/{sandboxID}
```

**响应 (200)：** 成功时为空。

---

## 错误格式

所有 API 错误均采用如下 JSON 结构：

```json
{
  "code": 400,
  "message": "可读的错误描述"
}
```

常见 HTTP 状态码：

| 状态码 | 含义 |
|--------|------|
| 400 | 请求错误（输入无效、参数缺失） |
| 401 | 未授权（缺少或无效的 API Key） |
| 404 | 沙箱未找到 |
| 500 | 内部服务错误 |

---

## 其他端点

| 路径 | 说明 |
|------|------|
| `/mcp` | MCP 服务端点（Streamable HTTP） |
| `/sandbox/{name}/` | 沙箱容器路由 |
| `/sandboxes/router/{sandboxID}/{port}/` | E2B 风格沙箱端口路由 |
| `/healthz` | 健康检查，返回 `{"status":"ok","version":"..."}` |
| `/ui/` | Web UI 静态文件 |