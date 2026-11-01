package app

import (
	"fmt"
	"os"
	"path/filepath"

	"my-app/internal/icons"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type pageImageData struct {
	Image       *ui.Bitmap
	Width       float32
	Height      float32
	Note        string
	initialized bool
}

func (app *app) pageImage(c *ui.Context) {
	d := &app.pageImageData
	if !d.initialized {
		d.initialized = true
		app.loadImage()
	}
	t := c.Theme()

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Image")
	})

	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "位图").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					if d.Image == nil {
						ui.Text(c, d.Note).TextColor(t.TextMuted)
						return
					}
					ui.Text(c, fmt.Sprintf("ui.DecodeBitmap 解码 icon.png，原始尺寸 %.0f × %.0f", d.Width, d.Height))
				})

				ui.Text(c, "ui.Image(c, bitmap)")
				if d.Image != nil {
					ui.Image(c, d.Image).Size(160, 160).Radius(8).
						Border(1, t.Border).JustifySelf(ui.Start)
				}
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "缩放模式").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "同一张图片放进 120 × 90 的框，用 Fit 控制缩放与裁切。")
				})

				ui.Text(c, "ui.Contain / Cover / FillBox / ScaleDown / NaturalSize")
				ui.Row(c).Gap(12).Wrap().Children(func() {
					for _, mode := range []struct {
						name string
						f    ui.Fit
					}{
						{"Contain", ui.Contain},
						{"Cover", ui.Cover},
						{"FillBox", ui.FillBox},
						{"ScaleDown", ui.ScaleDown},
						{"NaturalSize", ui.NaturalSize},
					} {
						ui.Column(c).Gap(6).Children(func() {
							box := ui.Box(c).Size(120, 90).Radius(6).
								Background(t.Surface).Border(1, t.Border)
							box.Children(func() {
								if d.Image != nil {
									ui.Image(c, d.Image).Fill().Fit(mode.f)
								}
							})
							ui.Text(c, mode.name).FontSize(12).TextColor(t.TextMuted)
						})
					}
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "圆角与圆形").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "Radius 控制圆角；半径不小于一半边长时即为圆形。")
				})

				ui.Text(c, "ui.Element.Radius")
				ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
					if d.Image != nil {
						ui.Image(c, d.Image).Size(120, 120).Radius(12).Border(1, t.Border)
						ui.Image(c, d.Image).Size(120, 120).Radius(60).Border(1, t.Border)
					}
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "SVG 图片").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "icons 里的图标是 SVG，可直接作为 Image 的图片源；currentColor 用文字颜色。")
				})

				ui.Text(c, "ui.Image(c, svg)")
				ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
					ui.Image(c, icons.Upload).Size(40, 40).TextColor(t.Accent)
					ui.Image(c, icons.Picture).Size(40, 40).TextColor(t.TextMuted)
					ui.Image(c, icons.Music).Size(40, 40).TextColor(ui.Hex("#f59e0b"))
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "支持的格式").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "ui.DecodeBitmap 支持 PNG、JPEG、GIF、WebP、BMP；ui.DecodeImage 还会尝试按 SVG 解析。")
					ui.Text(c, "图片只需解码一次，之后在视图里复用同一个 Bitmap。")
				})
			})
	})
}

// loadImage 从应用资源里读取 icon.png，作为图片示例的位图
func (app *app) loadImage() {
	d := &app.pageImageData
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
	d.Width, d.Height = ui.ImageSize(bmp)
	d.Image = bmp
}
