package app

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

type pageTreeData struct {
	initialized bool

	SrcOpen   bool
	TestsOpen bool
	Chosen    string
}

func (app *app) pageTree(c *ui.Context) {
	d := &app.pageTreeData
	if !d.initialized {
		d.SrcOpen = true
		d.initialized = true
	}

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Tree")
	})

	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Tree").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("已选择: %s", orNone(d.Chosen)))
				})

				ui.Text(c, "ui.Tree")
				ui.Tree(c, func() {
					if ui.TreeItem(c, "src", &d.SrcOpen, func() {
						if ui.TreeItem(c, "main.go", nil, nil).Clicked() {
							d.Chosen = "src/main.go"
						}
						if ui.TreeItem(c, "ui_layout.go", nil, nil).Clicked() {
							d.Chosen = "src/ui_layout.go"
						}
					}).Clicked() {
						d.Chosen = "src"
					}
					if ui.TreeItem(c, "tests", &d.TestsOpen, func() {
						if ui.TreeItem(c, "app_test.go", nil, nil).Clicked() {
							d.Chosen = "tests/app_test.go"
						}
					}).Clicked() {
						d.Chosen = "tests"
					}
					if ui.TreeItem(c, "go.mod", nil, nil).Clicked() {
						d.Chosen = "go.mod"
					}
				}).Height(240)
			})

		ui.Text(c, "点击箭头展开或折叠节点，点击条目或按回车可选中。").
			TextColor(c.Theme().TextMuted).Margin(PageTitleBarMargin)
	})
}
