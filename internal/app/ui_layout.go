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

func (app *app) view(c *ui.Context) {
	app.setTitlebar(c) // 自定义标题栏
	// 内容区域
	ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
		app.setSidebar(c)
		app.setRouter(c)
	})
}

// setTitlebar 设置标题栏
func (app *app) setTitlebar(c *ui.Context) {
	bar := c.TitleBar()
	ui.Row(c).Height(max(bar.Height, 32)).Padding(0, bar.Right+12, 0, bar.Left+12).AlignItems(ui.Center).DragWindow().Children(func() {
		ui.Icon(c, icons.Home).Size(16, 16)
		ui.Text(c, mygo.App.Name()).Bold().Margin(0, 0, 0, 8) // 左边留 8px 与图标间隔
		ui.Box(c).Grow(1)                                     // 占满剩余空间，把按钮推到最右侧
		app.setDarkMode(c)                                    // 暗黑模式切换
	})
}

// setDarkMode 暗黑与普通模式间切换
func (app *app) setDarkMode(c *ui.Context) {
	t := c.Theme()
	label, icon := "切换到暗黑模式", icons.Moon
	if mygo.Theme.IsDark() {
		label, icon = "切换到普通模式", icons.Sun
	}
	button := ui.ButtonBase(c.Key("theme-toggle")).Size(32, 28).Radius(6).Center().
		TextColor(t.TextMuted).Label(label).Tooltip(label)
	if button.Hovered() {
		button.Background(t.SurfaceHover)
	}
	button.Children(func() { ui.Icon(c, icon).Size(16, 16) })
	if button.Clicked() {
		if mygo.Theme.IsDark() {
			mygo.Theme.SetSource(mygo.ThemeLight)
		} else {
			mygo.Theme.SetSource(mygo.ThemeDark)
		}
	}
}

// setSidebar 设置侧边栏
func (app *app) setSidebar(c *ui.Context) {
	page := app.router.Path()
	ui.Sidebar(c, &page, func() {
		ui.SidebarItem(c, "/home", icons.Home, "首页")
		ui.SidebarItem(c, "/text", icons.Save, "Text")
		ui.SidebarItem(c, "/input", icons.Check, "Input")
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
		case r.Match("/input"):
			app.pageInput(c)
		default:
			ui.Text(c, "Not found")
		}
	})
}
