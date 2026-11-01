package app

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type pageClipboardData struct {
	Text        string
	ClipText    string
	ClipHTML    string
	Formats     string
	Image       ui.ImageSource
	ImageNote   string
	initialized bool
}

func (app *app) pageClipboard(c *ui.Context) {
	d := &app.pageClipboardData
	if !d.initialized {
		d.Text = "Hello, mygo!"
		d.initialized = true
	}

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Clipboard")
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Text").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("剪贴板文本: %s", orNone(d.ClipText)))
				})

				ui.Text(c, "mygo.Clipboard.ReadText / WriteText")
				ui.TextInput(c, &d.Text).Label("文本")
				ui.Row(c).Gap(8).Children(func() {
					if ui.Button(c, "写入").Clicked() {
						mygo.Clipboard.WriteText(d.Text)
					}
					if ui.Button(c, "读取").Clicked() {
						d.ClipText = mygo.Clipboard.ReadText()
					}
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "HTML").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("剪贴板 HTML: %s", orNone(d.ClipHTML)))
				})

				ui.Text(c, "mygo.Clipboard.ReadHTML / WriteHTML")
				ui.Row(c).Gap(8).Children(func() {
					if ui.Button(c, "写入 HTML").Clicked() {
						mygo.Clipboard.WriteHTML("<b>mygo</b> 剪贴板示例")
					}
					if ui.Button(c, "读取 HTML").Clicked() {
						d.ClipHTML = mygo.Clipboard.ReadHTML()
					}
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Image").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, orNone(d.ImageNote))
					if d.Image != nil {
						ui.Image(c, d.Image).Size(64, 64)
					}
				})

				ui.Text(c, "mygo.Clipboard.ReadImage")
				ui.Row(c).Gap(8).Children(func() {
					if ui.Button(c, "读取图片").Clicked() {
						app.readClipboardImage()
					}
					if ui.Button(c, "清空剪贴板").Clicked() {
						mygo.Clipboard.Clear()
						d.Image = nil
						d.ImageNote = "已清空剪贴板"
					}
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Formats").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, orNone(d.Formats))
				})

				ui.Text(c, "mygo.Clipboard.AvailableFormats")
				if ui.Button(c, "列出格式").JustifySelf(ui.Start).Clicked() {
					formats := mygo.Clipboard.AvailableFormats()
					d.Formats = strings.Join(formats, "、")
				}
			})
	})
}

// readClipboardImage 读取剪贴板中的图片并解码为可显示的图像
func (app *app) readClipboardImage() {
	d := &app.pageClipboardData
	png := mygo.Clipboard.ReadImage()
	if len(png) == 0 {
		d.Image = nil
		d.ImageNote = "剪贴板中没有图片"
		return
	}
	src, err := ui.DecodeImage(png)
	if err != nil {
		d.Image = nil
		d.ImageNote = fmt.Sprintf("无法解析图片: %v", err)
		return
	}
	d.Image = src
	d.ImageNote = fmt.Sprintf("PNG，%d 字节", len(png))
}
