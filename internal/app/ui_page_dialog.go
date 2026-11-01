package app

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type pageDialogData struct {
	Files     []string // 最近一次选择的文件
	Directory string   // 最近一次选择的目录
	SavePath  string   // 最近一次选择的保存路径
	Message   string   // 最近一次消息对话框的结果
	Error     string   // 最近一次错误对话框的结果
}

func (app *app) pageDialog(c *ui.Context) {
	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Dialog")
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "OpenFile").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("已选择: %s", orNone(strings.Join(app.pageDialogData.Files, "、"))))
				})

				ui.Text(c, "mygo.Dialog.Open")
				if ui.Button(c, "选择文件…").JustifySelf(ui.Start).Clicked() {
					go app.pickFiles()
				}
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "OpenDir").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("已选择: %s", orNone(app.pageDialogData.Directory)))
				})

				ui.Text(c, "mygo.Dialog.Open (Directory)")
				if ui.Button(c, "打开目录…").JustifySelf(ui.Start).Clicked() {
					go app.pickDirectory()
				}
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "SaveFile").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("保存到: %s", orNone(app.pageDialogData.SavePath)))
				})

				ui.Text(c, "mygo.Dialog.Save")
				if ui.Button(c, "保存文件…").JustifySelf(ui.Start).Clicked() {
					go app.saveFile()
				}
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Message").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("结果: %s", orNone(app.pageDialogData.Message)))
				})

				ui.Text(c, "mygo.Dialog.Message")
				if ui.Button(c, "显示消息…").JustifySelf(ui.Start).Clicked() {
					go app.showMessage()
				}
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Error").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("结果: %s", orNone(app.pageDialogData.Error)))
				})

				ui.Text(c, "mygo.Dialog.Message (MessageError)")
				if ui.Button(c, "显示错误…").JustifySelf(ui.Start).Clicked() {
					go app.showError()
				}
			})
	})
}

// pickFiles 在后台打开系统的文件选择框，选择结果交回主线程写入页面状态
func (app *app) pickFiles() {
	paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
		Parent:   app.window,
		Title:    "选择文件",
		Multiple: true,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("选择文件: %v", err))
		return
	}
	app.window.Update(func() { app.pageDialogData.Files = paths })
}

// pickDirectory 在后台打开系统的目录选择框
func (app *app) pickDirectory() {
	paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
		Parent:    app.window,
		Title:     "打开目录",
		Directory: true,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("打开目录: %v", err))
		return
	}
	var dir string
	if len(paths) > 0 {
		dir = paths[0]
	}
	app.window.Update(func() { app.pageDialogData.Directory = dir })
}

// saveFile 在后台打开系统的保存对话框
func (app *app) saveFile() {
	path, err := mygo.Dialog.Save(mygo.SaveDialogOptions{
		Parent:      app.window,
		Title:       "保存文件",
		DefaultPath: "notes.txt",
		Filters:     []mygo.FileFilter{{Name: "文本", Extensions: []string{"txt", "md"}}},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("保存文件: %v", err))
		return
	}
	app.window.Update(func() { app.pageDialogData.SavePath = path })
}

// showMessage 在后台弹出系统的消息对话框，结果交回主线程
func (app *app) showMessage() {
	res, err := mygo.Dialog.Message(mygo.MessageOptions{
		Parent:  app.window,
		Type:    mygo.MessageInfo,
		Title:   "提示",
		Message: "操作已完成。",
		Detail:  "这里是放在小字里的详细说明。",
	})
	if err != nil {
		slog.Error(fmt.Sprintf("消息对话框: %v", err))
		return
	}
	result := fmt.Sprintf("点击了第 %d 个按钮", res.Button+1)
	app.window.Update(func() { app.pageDialogData.Message = result })
}

// showError 在后台弹出系统的错误对话框，并区分“重试”与“取消”
func (app *app) showError() {
	res, err := mygo.Dialog.Message(mygo.MessageOptions{
		Parent:       app.window,
		Type:         mygo.MessageError,
		Title:        "出错了",
		Message:      "无法写入文件。",
		Detail:       "请检查磁盘空间与文件权限后重试。",
		Buttons:      []string{"重试", "取消"},
		CancelButton: 1,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("错误对话框: %v", err))
		return
	}
	action := "取消"
	if res.Button == 0 {
		action = "重试"
	}
	app.window.Update(func() { app.pageDialogData.Error = action })
}

// orNone 在未选择内容时返回占位符
func orNone(s string) string {
	if s == "" {
		return "（无）"
	}
	return s
}
