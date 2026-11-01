package app

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type pageDialogData struct {
	DialogOpen bool
	AlertOpen  bool
	Name       string
	LastAction string
	Files      []string // 最近一次选择的文件
	Directory  string   // 最近一次选择的目录
	SavePath   string   // 最近一次选择的保存路径
}

func (app *app) pageDialog(c *ui.Context) {
	ui.Row(c).Children(func() {
		ui.Text(c, "Dialog").FontSize(PageTitleFontSize).Margin(PageTitleBarMargin).Bold()
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Dialog").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("对话框打开: %t", app.pageDialogData.DialogOpen))
					ui.Text(c, fmt.Sprintf("最近操作: %s", app.pageDialogData.LastAction))
				})

				ui.Text(c, "ui.Modal")
				if ui.Button(c, "打开对话框").Clicked() {
					app.pageDialogData.DialogOpen = true
				}
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "AlertDialog").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("确认框打开: %t", app.pageDialogData.AlertOpen))
				})

				ui.Text(c, "ui.AlertDialog")
				if ui.Button(c, "删除…").Clicked() {
					app.pageDialogData.AlertOpen = true
				}
			})

		ui.Divider(c)

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
	})

	// 对话框是覆盖层，放在页面内容之后构建
	ui.Modal(c, &app.pageDialogData.DialogOpen, func() {
		ui.Text(c, "重命名").FontSize(18).Bold()
		ui.Text(c, "输入新的名称。").TextColor(c.Theme().TextMuted)
		ui.TextInput(c, &app.pageDialogData.Name).AutoFocus()
		ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
			if ui.Button(c, "取消").Clicked() {
				app.pageDialogData.DialogOpen = false
			}
			if ui.PrimaryButton(c, "确定").Clicked() {
				app.pageDialogData.LastAction = fmt.Sprintf("重命名为 %q", app.pageDialogData.Name)
				app.pageDialogData.DialogOpen = false
			}
		})
	})

	switch ui.AlertDialog(c, &app.pageDialogData.AlertOpen, "删除这条记录？", "删除后无法恢复。", "取消", "删除") {
	case ui.AlertDismissed:
		app.pageDialogData.LastAction = "已取消删除"
	case 0:
		app.pageDialogData.LastAction = "已取消删除"
	case 1:
		app.pageDialogData.LastAction = "已删除"
	}
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

// orNone 在未选择内容时返回占位符
func orNone(s string) string {
	if s == "" {
		return "（无）"
	}
	return s
}
