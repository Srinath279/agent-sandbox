---
icon: lucide/rocket
---

# Agent-Sandbox 概述

![Agent-Sandbox](assets/light.jpg#only-light)
![Agent-Sandbox](assets/dark.jpg#only-dark)


Agent-Sandbox 是一个开源的、Kubernetes 原生的 AI Agent 运行时平台。
它为代码执行、浏览器/计算机任务和 Shell 工作流提供隔离的、有状态的、多租户的沙箱环境。

该项目设计兼容 E2B 协议和 SDK 工作流，同时提供原生 REST API、MCP 和 Web UI。

---

## 特性

<div class="grid cards" markdown>

-   :material-clock-fast:{ .lg .middle } __5 分钟快速部署__

    ---

    使用 [`kubectl apply`](#) 安装 [`Agent-Sandbox`](#)，几分钟即可启动运行

    [:octicons-arrow-right-24: 快速开始](quickstart)

-   :lucide-plug:{ .lg .middle } __AI 友好__

    ---

    支持 Skills、CLI 和 MCP，为你的 Agent 赋能沙箱能力

    [:octicons-arrow-right-24: 参考](cli)

-   :lucide-mouse-pointer-click:{ .lg .middle } __简单易用__

    ---

    设计简洁的 **REST API** 和 **UI**，最小化 Kubernetes 对象部署，易于使用和维护。

-   :lucide-package-open:{ .lg .middle } __开放灵活__

    ---
   
    支持 E2B 模板和任意容器镜像，可通过 **自定义模版** 扩展。

</div>

---

## 核心能力

- **E2B 协议兼容**：支持沙箱生命周期 API 和路由。
- **沙箱生命周期管理**：创建、列表、连接、删除。
- **多租户访问控制**：支持系统用户和普通用户。
- **模版和池管理**：实现快速沙箱分配和预热。
- **可观测性**：沙箱事件、指标和日志。
- **交互式操作**：终端、文件上传/下载、沙箱路由。
- **MCP 服务集成**：支持 Agent 原生自动化。
- **Web UI**：沙箱/模版/池操作和运行时检查界面。

## 快速开始 :octicons-heart-fill-24:{ .heart }

参见 [快速开始](quickstart.md)，了解部署、通过 E2B SDK 创建沙箱以及使用 REST API的分步指南。