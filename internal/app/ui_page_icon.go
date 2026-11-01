package app

import (
	"my-app/internal/icons"

	"github.com/egoist/mygo/ui"
)

// pageIcons 是 icons 包里可供示例展示的全部图标
var pageIcons = []struct {
	Name string
	Icon *ui.SVG
}{
	{"Send", icons.Send},
	{"Search", icons.Search},
	{"Chevron", icons.Chevron},
	{"Plus", icons.Plus},
	{"Close", icons.Close},
	{"Folder", icons.Folder},
	{"Sliders", icons.Sliders},
	{"Sun", icons.Sun},
	{"Gear", icons.Gear},
	{"Monitor", icons.Monitor},
	{"Moon", icons.Moon},
	{"MoonLit", icons.MoonLit},
	{"Translate", icons.Translate},
	{"Braces", icons.Braces},
	{"Triangle", icons.Triangle},
	{"Radio", icons.Radio},
	{"Copy", icons.Copy},
	{"Trash", icons.Trash},
	{"Grip", icons.Grip},
	{"Alert", icons.Alert},
	{"Check", icons.Check},
	{"Eye", icons.Eye},
	{"Home", icons.Home},
	{"More", icons.More},
	{"MoreVertical", icons.MoreVertical},
	{"Plug", icons.Plug},
	{"SidebarClose", icons.SidebarClose},
	{"SidebarOpen", icons.SidebarOpen},
	{"Bolt", icons.Bolt},
	{"Stack", icons.Stack},
	{"External", icons.External},
	{"Pencil", icons.Pencil},
	{"Wand", icons.Wand},
	{"File", icons.File},
	{"Play", icons.Play},
	{"Pause", icons.Pause},
	{"Music", icons.Music},
	{"Download", icons.Download},
	{"Picture", icons.Picture},
	{"Video", icons.Video},
	{"Save", icons.Save},
	{"Upload", icons.Upload},
	{"Doc", icons.Doc},
	{"Type", icons.Type},
	{"TextCursor", icons.TextCursor},
	{"CheckSquare", icons.CheckSquare},
	{"List", icons.List},
	{"Table", icons.Table},
	{"Tree", icons.Tree},
	{"Grid", icons.Grid},
	{"Code", icons.Code},
	{"Activity", icons.Activity},
	{"Message", icons.Message},
	{"Dialog", icons.Dialog},
	{"Bell", icons.Bell},
	{"Clipboard", icons.Clipboard},
	{"Terminal", icons.Terminal},
	{"Keyboard", icons.Keyboard},
	{"Chart", icons.Chart},
	{"User", icons.User},
}

func (app *app) pageIcon(c *ui.Context) {
	t := c.Theme()

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Icon")
	})

	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "尺寸").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "图标默认与字号同高，可用 Size 指定宽高，保持 SVG 比例。")
				})

				ui.Text(c, "ui.Icon(c, svg).Size(w, h)")
				ui.Row(c).Gap(24).AlignItems(ui.Center).Children(func() {
					for _, size := range []float32{16, 24, 32, 48} {
						ui.Column(c).Gap(6).AlignItems(ui.Center).Children(func() {
							ui.Icon(c, icons.Sun).Size(size, size).TextColor(t.Accent)
							ui.Textf(c, "%.0f", size).FontSize(11).TextColor(t.TextMuted)
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
				ui.Text(c, "颜色").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "描边图标使用 currentColor，因此跟随 TextColor。")
				})

				ui.Text(c, "TextColor")
				ui.Row(c).Gap(16).Children(func() {
					ui.Icon(c, icons.Bolt).Size(28, 28).TextColor(t.Accent)
					ui.Icon(c, icons.Bolt).Size(28, 28).TextColor(t.TextMuted)
					ui.Icon(c, icons.Bolt).Size(28, 28).TextColor(ui.Hex("#ef4444"))
					ui.Icon(c, icons.Bolt).Size(28, 28).TextColor(ui.Hex("#10b981"))
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "与其他元素组合").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "图标是装饰元素，默认对辅助技术不可见；需要时可加 Label。")
				})

				ui.Text(c, "图标 + 文字 / 按钮")
				ui.Row(c).Gap(20).AlignItems(ui.Center).Children(func() {
					ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
						ui.Icon(c, icons.Folder).Size(16, 16)
						ui.Text(c, "项目文件")
					})

					ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
						ui.Icon(c, icons.Search).Size(16, 16).TextColor(t.TextMuted)
						ui.Text(c, "搜索").TextColor(t.TextMuted)
					})

					button := ui.ButtonBase(c).Padding(6, 12).Radius(6).
						Background(t.Accent).TextColor(t.AccentText).Label("保存")
					button.Children(func() {
						ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
							ui.Icon(c, icons.Save).Size(16, 16)
							ui.Text(c, "保存")
						})
					})
					if button.Hovered() {
						button.Background(t.Accent.Alpha(0.85))
					}
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "全部图标").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Textf(c, "icons 包共有 %d 个图标，均以 24 × 24 描边风格定义。", len(pageIcons)).
						TextColor(t.TextMuted)
				})

				ui.Text(c, "icons 包")
				ui.Row(c).Wrap().Gap(12).Children(func() {
					for _, item := range pageIcons {
						ui.Column(c).Width(84).Padding(8).Radius(6).Gap(6).
							Background(t.Surface).AlignItems(ui.Center).Children(func() {
							ui.Icon(c, item.Icon).Size(24, 24)
							ui.Text(c, item.Name).FontSize(11).TextColor(t.TextMuted).
								MaxLines(1).Tooltip(item.Name)
						})
					}
				})
			})
	})
}
