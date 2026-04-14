---
icon: lucide/rocket
---

# Agent-Sandbox 概览

Agent-Sandbox 是一个面向 AI Agent 的开源沙箱运行时平台，基于 Kubernetes，提供隔离、可持久、多租户的沙箱能力，适用于代码执行、浏览器/桌面操作、终端命令等场景。

项目支持 E2B 协议与 SDK 生态，同时提供原生 REST API、Skills、MCP 接口和 Web UI。

## 核心能力

- **兼容 E2B 协议**，可接入现有 E2B 风格工作流。
- **沙箱全生命周期管理**：创建、查询、连接、删除。
- **多租户与权限控制**：区分系统用户与普通用户。
- **模板与池化能力**：支持模板池、预热和动态模板。
- **可观测性**：支持事件、指标、日志。
- **交互式运维能力**：终端、文件上传下载、沙箱路由访问。
- **MCP 集成**：可供 Agent 自动化管理沙箱。
- **可视化 UI**：统一管理沙箱、模板、池和运行状态。

## 接口与入口

Agent-Sandbox 提供两类 API 接口：

- **原生 REST API**（`/api/v1`）：完整的沙箱生命周期、终端、文件、配置、池、日志、事件、指标。
- **E2B 兼容 API**（`/e2b/v1`）：与 E2B 协议兼容的沙箱 CRUD 和连接接口。

此外还提供 MCP（`/mcp`）、沙箱路由、UI 和健康检查端点。

详见 [API 参考](api.md)。

## 鉴权与权限

- HTTP 层已启用 API Key 鉴权中间件。
- 支持系统用户与普通用户权限隔离。
- 文档中的默认系统用户 token：
  - `sys-2492a85b10ed4cb083b2c76b181eac96`
- 可通过环境变量配置用户 token。

## 快速开始

详见 [快速开始](quickstart.md)，包含部署、通过 E2B SDK 创建沙箱和使用 REST API 的完整步骤。

## E2B CLI 兼容说明

CLI 使用说明见：

- `/cli/README.md`

其中包含 Installation、Authentication、List sandboxes、Create sandbox、Connect to sandbox、Execute commands in sandbox 等章节。

## 版本演进摘要

来自 `CHANGELOG.md` 的关键版本：

- **v0.2.0**：支持 E2B 协议、Code Interpreter、Desktop。
- **v0.3.0**：引入模板池和动态模板。
- **v0.3.1**：支持从 ConfigMap 动态加载模板配置、模板池预热。
- **v0.4.0**：上线沙箱管理 UI。
- **v0.4.1**：增加事件/指标能力与模板资源、元数据配置。
- **v0.4.2**：增加模板 args、基础权限控制、池获取性能优化、UI 排序增强。

## 参考

- 根目录 README：`/README.md`
- 变更日志：`/CHANGELOG.md`
- 仓库内 E2B CLI 指南：`/cli/README.md`
