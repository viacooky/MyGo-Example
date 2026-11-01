package app

import (
	"fmt"
	"path/filepath"

	"my-app/internal/icons"

	"github.com/egoist/mygo/ui"
)

type pageOutlineData struct {
	Outline  ui.OutlineState[string]
	Chosen   ui.Selection[string]
	Children map[string][]string
}

func (app *app) pageOutline(c *ui.Context) {
	d := &app.pageOutlineData
	if d.Children == nil {
		d.Children = map[string][]string{
			"src":     {"src/app", "src/ui"},
			"src/app": {"src/app/main.go", "src/app/layout.go"},
			"src/ui":  {"src/ui/button.go", "src/ui/list.go"},
		}
		d.Outline.Open.Add("src")
		d.Outline.Open.Add("src/app")
	}
	d.Outline.List.Selection = &d.Chosen

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Outline")
	})

	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Outline").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("已选择 %d 项（支持 Ctrl / Shift 多选）", d.Chosen.Len()))
				})

				ui.Text(c, "ui.Outline")
				ui.Outline(c, &d.Outline, []string{"src"}, func(item string) []string {
					return d.Children[item]
				}, func(item string) {
					icon := icons.File
					if len(d.Children[item]) > 0 {
						icon = icons.Folder
					}
					ui.Row(c).Gap(6).Children(func() {
						ui.Icon(c, icon).Size(14, 14)
						ui.Text(c, filepath.Base(item))
					})
				}).Height(280)
			})

		ui.Text(c, "点击箭头展开或折叠，点击条目选中；目录与文件由 children 回调决定。").
			TextColor(c.Theme().TextMuted).Margin(PageTitleBarMargin)
	})
}
