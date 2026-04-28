---
icon: lucide/rocket
---

# 快速开始

完整文档请访问 [zensical.org](https://zensical.org/docs/)。

## 命令

* [`zensical new`][new] - 创建新项目
* [`zensical serve`][serve] - 启动本地 Web 服务
* [`zensical build`][build] - 构建站点

  [new]: https://zensical.org/docs/usage/new/
  [serve]: https://zensical.org/docs/usage/preview/
  [build]: https://zensical.org/docs/usage/build/

## 示例

### 提示框

> 参见 [文档](https://zensical.org/docs/authoring/admonitions/)

!!! note

    这是一个 **note（注意）** 提示框。用于提供有用的信息。

!!! warning

    这是一个 **warning（警告）** 提示框。请注意！

### 折叠块

> 参见 [文档](https://zensical.org/docs/authoring/admonitions/#collapsible-blocks)

??? info "点击展开查看更多信息"

    这段内容默认隐藏，点击后展开。
    适合 FAQ 或较长的说明。

## 代码块

> 参见 [文档](https://zensical.org/docs/authoring/code-blocks/)

``` python hl_lines="2" title="代码块"
def greet(name):
    print(f"Hello, {name}!") # (1)!

greet("Python")
```

1.  > 参见 [文档](https://zensical.org/docs/authoring/code-blocks/#code-annotations)

    代码注释允许为代码行附加说明。

代码也可以内联高亮：`#!python print("Hello, Python!")`。

## 内容标签页

> 参见 [文档](https://zensical.org/docs/authoring/content-tabs/)

=== "Python"

    ``` python
    print("Hello from Python!")
    ```

=== "Rust"

    ``` rs
    println!("Hello from Rust!");
    ```

## 图表

> 参见 [文档](https://zensical.org/docs/authoring/diagrams/)

``` mermaid
graph LR
  A[Start] --> B{Error?};
  B -->|Yes| C[Hmm...];
  C --> D[Debug];
  D --> B;
  B ---->|No| E[Yay!];
```

## 脚注

> 参见 [文档](https://zensical.org/docs/authoring/footnotes/)

这是一句带脚注的句子。[^1]

悬停查看提示信息。

[^1]: 这是脚注内容。


## 格式化

> 参见 [文档](https://zensical.org/docs/authoring/formatting/)

- ==这是标记（高亮）==
- ^^这是插入（下划线）^^
- ~~这是删除（删除线）~~
- H~2~O
- A^T^A
- ++ctrl+alt+del++

## 图标和表情

> 参见 [文档](https://zensical.org/docs/authoring/icons-emojis/)

* :sparkles: `:sparkles:`
* :rocket: `:rocket:`
* :tada: `:tada:`
* :memo: `:memo:`
* :eyes: `:eyes:`

## 数学公式

> 参见 [文档](https://zensical.org/docs/authoring/math/)

$$
\cos x=\sum_{k=0}^{\infty}\frac{(-1)^k}{(2k)!}x^{2k}
$$

!!! warning "需要配置"
    注意 MathJax 是通过本页的 `script` 标签引入的，并没有在生成的默认配置中启用，
    以避免在不需要数学公式的页面中加载。如果你的页面数学公式较多，请参考文档了解如何在所有页面配置它。

<script id="MathJax-script" src="https://unpkg.com/mathjax@3/es5/tex-mml-chtml.js"></script>
<script>
  window.MathJax = {
    tex: {
      inlineMath: [["\\(", "\\)"]],
      displayMath: [["\\[", "\\]"]],
      processEscapes: true,
      processEnvironments: true
    },
    options: {
      ignoreHtmlClass: ".*|",
      processHtmlClass: "arithmatex"
    }
  };

  document$.subscribe(() => {
    MathJax.startup.output.clearCache()
    MathJax.typesetClear()
    MathJax.texReset()
    MathJax.typesetPromise()
  })
</script>

## 任务列表

> 参见 [文档](https://zensical.org/docs/authoring/lists/#using-task-lists)

* [x] 安装 Zensical
* [x] 配置 `zensical.toml`
* [x] 编写精彩文档
* [ ] 部署到任意位置

## 工具提示

> 参见 [文档](https://zensical.org/docs/authoring/tooltips/)

[悬停我][example]

  [example]: https://example.com "我是一个工具提示！"