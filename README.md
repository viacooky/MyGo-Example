# MyGo-Example

MyGo-Example 是一个基于 [MyGo](https://github.com/egoist/mygo) 的示例项目，用于展示该框架的基本用法：原生 UI、路由式多页面布局、自绘标题栏、托盘与输入控件。

## 快速开始

环境要求：

- Go 1.27.3 及以上（版本见 `go.mod`）
- 本项目使用原生 UI，不需要 Bun、Node 或其他前端工具链
- Linux 下托盘图标依赖 `libayatana-appindicator3`；缺失时应用仍可运行，只是没有托盘图标

克隆并启动：

```sh
git clone <repo-url> mygo-lean
cd mygo-lean
go tool mygo dev
```

MyGo 作为 Go 工具依赖登记在 `go.mod` 的 `tool` 指令中，首次执行会自动下载，无需全局安装。

## 开发与构建 (Development / Building)

### 开发

开发模式带热重载，修改代码后窗口会重建：

```sh
go tool mygo dev
```

也可以直接运行入口包：

```sh
go run ./cmd/gui
```

入口 `cmd/gui/main.go` 只负责调用 `internal/app.Run()`；窗口、路由、布局和各页面都在 `internal/app`，图标集中在 `internal/icons`（内嵌 SVG）。`mygo dev` 与 `go run` 的区别在于前者会按 `mygo.json` 注入开发环境信息并提供热重载。

提交前建议执行检查：

```sh
go build ./...
go vet ./...
go tool mygo vet .
```

其中 `go tool mygo vet .` 是 MyGo 专用检查，会覆盖 UI 生命周期与线程使用等框架相关规则。

调试使用 `.vscode/launch.json` 中的配置：`Debug my-app` 以 debug 模式启动 `cmd/gui`；`Attach to my-app` 附加到已运行的进程，注意 `mygo dev` 的热重载会重启进程并断开调试会话。

### 构建

```sh
go tool mygo build
```

产物输出到 `mygo.json` 中 `out` 指定的目录（当前为 `dist/`）。应用名称、标识符、版本同样取自该文件，打包与分发时无需再改代码。

发布前可再用 `go tool mygo doctor` 检查项目配置。

### 目录结构

```
mygo-lean/
├── cmd/
│   └── gui/
│       └── main.go        程序入口，调用 internal/app.Run
├── internal/
│   ├── app/               窗口、路由、布局与各页面
│   └── icons/             内嵌的 SVG 图标
├── resources/
│   └── icon.png           应用图标，随包分发
├── .vscode/
│   └── launch.json        调试配置
├── mygo.json              应用元数据与构建配置
├── go.mod
└── go.sum
```
