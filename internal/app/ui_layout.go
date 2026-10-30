package app

import (
	"my-app/internal/icons"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

const (
	PageTitleFontSize  = 20
	PageTitleBarMargin = 8
)

// setTitlebar 设置标题栏
func (app *app) setTitlebar(c *ui.Context) {
	bar := c.TitleBar()
	ui.Row(c).Height(max(bar.Height, 32)).Padding(0, bar.Right+12, 0, bar.Left+12).DragWindow().Children(func() {
		ui.Text(c, mygo.App.Name()).Bold()
	})
}

// setSidebar 设置侧边栏
func (app *app) setSidebar(c *ui.Context) {
	page := app.router.Path()
	ui.Sidebar(c, &page, func() {
		ui.SidebarItem(c, "/home", icons.Home, "首页")
		ui.SidebarItem(c, "/text", icons.Save, "文本")
		ui.SidebarItem(c, "/page2", icons.Check, "页面2")
	}).
		Width(200).
		OnChange(func() {
			app.router.Push(page)
		})
}

// setRouter 设置路由 对应sidebar的路径
func (app *app) setRouter(c *ui.Context) {
	app.router.View(c, func(r *ui.Route) {
		switch {
		case r.Match("/home"):
			app.pageHome(c)
		case r.Match("/text"):
			app.pageText(c)
		case r.Match("/page2"):
			r.Title("页面2")
			app.page2(c)
		default:
			ui.Text(c, "Not found")
		}
	})
}
