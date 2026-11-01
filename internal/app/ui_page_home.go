package app

import (
	"runtime/debug"

	"github.com/egoist/mygo/ui"
)

func (app *app) pageHome(c *ui.Context) {
	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "概览")
	})

	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Text(c, "页面导航").FontSize(16).Bold().Margin(PageTitleBarMargin)
		ui.Grid(c).Columns(3).Gap(12).Margin(PageTitleBarMargin).Children(func() {
			navigate := func(label, path string) {
				if ui.Button(c, label).Clicked() {
					app.router.Push(path)
				}
			}
			navigate("Text", "/text")
			navigate("Input", "/input")
			navigate("Choices", "/choices")
			navigate("Actions", "/actions")
			navigate("DragDrop", "/dragdrop")
			navigate("List", "/list")
			navigate("Table", "/table")
			navigate("Tree", "/tree")
			navigate("Outline", "/outline")
			navigate("GridView", "/gridview")
			navigate("Document", "/document")
			navigate("Code", "/code")
			navigate("Status", "/status")
			navigate("Message", "/message")
			navigate("Modal", "/modal")
			navigate("Dialog", "/dialog")
			navigate("Notification", "/notification")
			navigate("Clipboard", "/clipboard")
			navigate("Shell", "/shell")
			navigate("GlobalShortcut", "/shortcut")
			navigate("Chart", "/charts")
			navigate("Image", "/image")
			navigate("Icon", "/icon")
			navigate("Avatar", "/avatar")
			navigate("Video", "/video")
			navigate("Audio", "/audio")
			navigate("Terminal", "/terminal")
		})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "关于").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "MyGo 组件示例应用，演示各组件的用法。")
				})

				ui.Text(c, "mygo 版本")
				ui.Text(c, mygoVersion())

				ui.Text(c, "GitHub")
				ui.Link(c, "github.com/egoist/mygo", "https://github.com/egoist/mygo")

				ui.Text(c, "项目")
				ui.Link(c, "github.com/viacooky/mygo-lean", "https://github.com/viacooky/mygo-lean")
			})
	})
}

// mygoVersion 返回当前构建实际使用的 mygo 版本
func mygoVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "未知"
	}
	for _, dep := range info.Deps {
		if dep.Path == "github.com/egoist/mygo" {
			return dep.Version
		}
	}
	return "未知"
}
