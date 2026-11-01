package app

import (
	"fmt"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type pageShellData struct {
	URL         string
	Path        string
	Result      string
	Icon        ui.ImageSource
	IconNote    string
	initialized bool
}

func (app *app) pageShell(c *ui.Context) {
	d := &app.pageShellData
	if !d.initialized {
		d.URL = "https://github.com/egoist/mygo"
		d.Path = "."
		d.initialized = true
	}

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Shell")
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "OpenExternal").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("结果: %s", orNone(d.Result)))
				})

				ui.Text(c, "mygo.Shell.OpenExternal")
				ui.TextInput(c, &d.URL).Label("URL")
				if ui.Button(c, "打开链接").JustifySelf(ui.Start).Clicked() {
					if err := mygo.Shell.OpenExternal(d.URL); err != nil {
						d.Result = fmt.Sprintf("打开失败: %v", err)
					} else {
						d.Result = "已在默认浏览器中打开"
					}
				}
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "OpenPath / ShowItemInFolder").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "用默认应用打开文件或目录，或在文件管理器中定位")
				})

				ui.Text(c, "mygo.Shell.OpenPath / ShowItemInFolder")
				ui.TextInput(c, &d.Path).Label("路径")
				ui.Row(c).Gap(8).Children(func() {
					if ui.Button(c, "打开路径").Clicked() {
						if err := mygo.Shell.OpenPath(d.Path); err != nil {
							d.Result = fmt.Sprintf("打开失败: %v", err)
						} else {
							d.Result = "已用默认应用打开"
						}
					}
					if ui.Button(c, "在文件管理器中显示").Clicked() {
						mygo.Shell.ShowItemInFolder(d.Path)
					}
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "FileIcon").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, orNone(d.IconNote))
					if d.Icon != nil {
						ui.Image(c, d.Icon).Size(48, 48)
					}
				})

				ui.Text(c, "mygo.Shell.FileIcon")
				if ui.Button(c, "读取上面的图标").JustifySelf(ui.Start).Clicked() {
					app.readShellIcon()
				}
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Beep").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "播放系统提示音")
				})

				ui.Text(c, "mygo.Shell.Beep")
				if ui.Button(c, "播放").JustifySelf(ui.Start).Clicked() {
					mygo.Shell.Beep()
				}
			})
	})
}

// readShellIcon 读取路径对应的系统图标并解码为可显示的图像
func (app *app) readShellIcon() {
	d := &app.pageShellData
	png, err := mygo.Shell.FileIcon(d.Path, 48)
	if err != nil {
		d.Icon = nil
		d.IconNote = fmt.Sprintf("读取图标失败: %v", err)
		return
	}
	src, err := ui.DecodeImage(png)
	if err != nil {
		d.Icon = nil
		d.IconNote = fmt.Sprintf("无法解析图标: %v", err)
		return
	}
	d.Icon = src
	d.IconNote = fmt.Sprintf("48×48 PNG，%d 字节", len(png))
}
