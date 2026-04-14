---
icon: lucide/rocket
---

# Agent-Sandbox Overview

![Agent-Sandbox](assets/light.jpg#only-light)
![Agent-Sandbox](assets/dark.jpg#only-dark)


Agent-Sandbox is an open-source, Kubernetes-native runtime platform for AI agents.
It provides isolated, stateful, multi-tenant sandboxes for code execution, browser/computer tasks, and shell workflows.

The project is designed to be compatible with E2B protocol and SDK workflows while also exposing native REST API, MCP, and a web UI.

---

## Features

<div class="grid cards" markdown>

-   :material-clock-fast:{ .lg .middle } __Set up in 5 minutes__

    ---

    Install [`Agent-Sandobx`](#) with [`kubectl apply`](#) and get up
    and running in minutes

    [:octicons-arrow-right-24: Getting started](quickstart)

-   :lucide-plug:{ .lg .middle } __Ai friendly__

    ---

    Supports Skills, Cli and MCP to empower your agents with sandbox capabilities

    [:octicons-arrow-right-24: Reference](cli)

-   :lucide-mouse-pointer-click:{ .lg .middle } __Easy to use__

    ---

    Designed with a simple **REST API**, **UI** and minimizing Kubernetes's objects to deploy for easy to use and maintain.

-   :lucide-package-open:{ .lg .middle } __Open and flexible__

    ---
   
    Supports E2B templates and any container image, and can be extended with **custom templates**.

</div>

---

## What it provides

- **E2B protocol compatibility** for sandbox lifecycle APIs and routing.
- **Sandbox lifecycle management**: create, list, connect, delete.
- **Multi-tenant access control** with system and regular users.
- **Template and pool management** for fast sandbox allocation and warmup.
- **Observability** with sandbox events, metrics, and logs.
- **Interactive operations**: terminal, file upload/download, and sandbox routing.
- **MCP server integration** for agent-native automation.
- **Web UI** for sandbox/template/pool operations and runtime inspection.

## Quick start :octicons-heart-fill-24:{ .heart }

See [Quick Start](quickstart.md) for a step-by-step guide covering deployment, creating a sandbox via E2B SDK, and using the REST API.
