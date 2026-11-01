package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type pageAvatarData struct {
	Name        string
	Image       *ui.Bitmap
	Note        string
	initialized bool
}

func (app *app) pageAvatar(c *ui.Context) {
	d := &app.pageAvatarData
	if !d.initialized {
		d.initialized = true
		d.Name = "Ada Lovelace"
		app.loadAvatarImage()
	}
	t := c.Theme()

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Avatar")
	})

	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "文字头像").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "没有图片时，Avatar 显示名字的首字母，并从名字里挑一个背景色。修改名字，头像随之变化。")
				})

				ui.Text(c, "ui.TextInput / ui.Avatar(c, name, nil)")
				ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
					ui.TextInput(c, &d.Name).Label("名字").Placeholder("输入名字").Width(220)
					ui.Avatar(c, d.Name, nil).Size(44, 44)
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "尺寸").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "Avatar 默认高度为字号的两倍，可用 Size 指定任意尺寸。")
				})

				ui.Text(c, "ui.Element.Size(w, h)")
				ui.Row(c).Gap(16).AlignItems(ui.End).Children(func() {
					for _, size := range []float32{28, 40, 56, 72} {
						ui.Avatar(c, d.Name, nil).Size(size, size)
					}
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "图片头像").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					if d.Image == nil {
						ui.Text(c, d.Note).TextColor(t.TextMuted)
						return
					}
					ui.Text(c, "传入 ui.Bitmap 后，图片会被裁成圆形；这里是应用的 icon.png。")
				})

				ui.Text(c, "ui.Avatar(c, name, bitmap)")
				ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
					if d.Image != nil {
						ui.Avatar(c, "MyGo", d.Image).Size(40, 40)
						ui.Avatar(c, "MyGo", d.Image).Size(56, 56)
						ui.Avatar(c, "MyGo", d.Image).Size(72, 72)
					}
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "辅助功能").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "Avatar 会以 name 作为图片名称暴露给辅助技术。")
				})

				ui.Text(c, "说明")
				ui.Text(c, "通常与列表、评论、账号菜单一起使用。").TextColor(t.TextMuted)
			})
	})
}

// loadAvatarImage 从应用资源里读取 icon.png，作为图片头像的位图
func (app *app) loadAvatarImage() {
	d := &app.pageAvatarData
	dir, err := mygo.App.Path(mygo.PathResources)
	if err != nil {
		d.Note = fmt.Sprintf("无法定位资源目录: %v", err)
		return
	}
	data, err := os.ReadFile(filepath.Join(dir, "icon.png"))
	if err != nil {
		d.Note = fmt.Sprintf("无法读取 icon.png: %v", err)
		return
	}
	bmp, err := ui.DecodeBitmap(data)
	if err != nil {
		d.Note = fmt.Sprintf("无法解码 icon.png: %v", err)
		return
	}
	d.Image = bmp
}
