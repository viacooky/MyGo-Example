package app

import "github.com/egoist/mygo/ui"

func (app *app) pageHome(c *ui.Context) {
	ui.Row(c).Children(func() {
		ui.Text(c, "首页").FontSize(PageTitleFontSize).Margin(PageTitleBarMargin).Bold()
	})

	ui.Row(c).Children(func() {

	})
}
