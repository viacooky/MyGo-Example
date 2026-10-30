package app

import (
	"github.com/egoist/mygo/ui"
)

func (app *app) view(c *ui.Context) {
	app.setTitlebar(c) // 自定义标题栏

	// 内容区域
	ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
		app.setSidebar(c)
		app.setRouter(c)
	})
}
