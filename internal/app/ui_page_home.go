package app

import "github.com/egoist/mygo/ui"

func (app *app) pageHome(c *ui.Context) {
	ui.Row(c).Children(func() {
		ui.Text(c, "首页").FontSize(PageTitleFontSize).Margin(PageTitleBarMargin).Bold()
	})

	ui.Grid(c).Columns(2).Gap(12).Margin(PageTitleBarMargin).Children(func() {
		if ui.Button(c, "Text").Clicked() {
			app.router.Push("/text")
		}
		if ui.Button(c, "Input").Clicked() {
			app.router.Push("/input")
		}
	})
}
