package app

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

type pageGridViewData struct {
	initialized bool

	Grid     ui.GridState
	Selected int
}

var gridViewColors = []ui.Color{
	ui.Hex("#4C6EF5"),
	ui.Hex("#12B886"),
	ui.Hex("#F59F00"),
	ui.Hex("#E64980"),
	ui.Hex("#7048E8"),
	ui.Hex("#15AABF"),
	ui.Hex("#F76707"),
	ui.Hex("#37B24D"),
}

func (app *app) pageGridView(c *ui.Context) {
	d := &app.pageGridViewData
	if !d.initialized {
		d.Selected = -1
		d.initialized = true
	}

	d.Grid.Selected = &d.Selected
	d.Grid.Label = func(i int) string { return fmt.Sprintf("第 %d 项", i+1) }

	const itemCount = 24

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "GridView")
	})

	// 让 GridView 直接占据剩余高度，由它自己的滚动容器处理滚轮，
	// 避免固定高度或嵌套滚动容器导致无法滚动。
	ui.Column(c).Fill().Padding(0, 16, 24, 16).Children(func() {
		ui.Text(c, fmt.Sprintf("已选择: %s", app.selectedGridItem())).Margin(PageTitleBarMargin)
		ui.GridView(c, &d.Grid, itemCount, 140, 100, func(i int) {
			background := gridViewColors[i%len(gridViewColors)]
			textColor := ui.Hex("#FFFFFF")
			if d.Selected == i {
				background = c.Theme().Accent
				textColor = c.Theme().AccentText
			}
			ui.Box(c).Fill().Radius(8).Background(background).Center().Children(func() {
				ui.Text(c, fmt.Sprintf("%02d", i+1)).Bold().TextColor(textColor)
			})
		}).Grow(1)

		ui.Text(c, "网格随窗口宽度自动换列，点击卡片即可选中。").
			TextColor(c.Theme().TextMuted).Margin(PageTitleBarMargin)
	})
}

// selectedGridItem 返回当前选中卡片的序号
func (app *app) selectedGridItem() string {
	if app.pageGridViewData.Selected < 0 {
		return ""
	}
	return fmt.Sprintf("第 %d 项", app.pageGridViewData.Selected+1)
}
