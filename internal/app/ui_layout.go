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

// breadcrumbTitle 以面包屑展示页面标题：首页 / 当前页，点击“首页”返回首页
func (app *app) breadcrumbTitle(c *ui.Context, title string) {
	items := []string{"首页", title}

	if ui.Breadcrumbs(c, items, &app.breadcrumbChosen).
		FontSize(PageTitleFontSize).
		Margin(PageTitleBarMargin).
		Bold().
		Label("页面导航").
		Changed() {
		app.router.Push("/home")
	}
	app.breadcrumbChosen = len(items) - 1
}

func (app *app) view(c *ui.Context) {
	app.setTitlebar(c) // 自定义标题栏
	// 内容区域：占据标题栏之外的剩余高度，侧边栏与页面才会各自滚动而不是撑出窗口
	ui.Row(c).Grow(1).AlignItems(ui.Stretch).Children(func() {
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
		ui.SidebarSection(c, "常用", nil, func() {
			ui.SidebarItem(c, "/text", icons.Type, "Text")
			ui.SidebarItem(c, "/input", icons.TextCursor, "Input")
			ui.SidebarItem(c, "/choices", icons.CheckSquare, "Choices")
			ui.SidebarItem(c, "/actions", icons.Sliders, "Actions")
			ui.SidebarItem(c, "/dragdrop", icons.Upload, "DragDrop")
		})
		ui.SidebarSection(c, "数据展示", nil, func() {
			ui.SidebarItem(c, "/list", icons.List, "List")
			ui.SidebarItem(c, "/table", icons.Table, "Table")
			ui.SidebarItem(c, "/tree", icons.Tree, "Tree")
			ui.SidebarItem(c, "/outline", icons.Folder, "Outline")
			ui.SidebarItem(c, "/gridview", icons.Grid, "GridView")
			ui.SidebarItem(c, "/document", icons.Doc, "Document")
			ui.SidebarItem(c, "/code", icons.Code, "Code")
		})
		ui.SidebarSection(c, "图表", nil, func() {
			ui.SidebarItem(c, "/charts", icons.Chart, "Chart")
		})
		ui.SidebarSection(c, "媒体", nil, func() {
			ui.SidebarItem(c, "/image", icons.Picture, "Image")
			ui.SidebarItem(c, "/icon", icons.Wand, "Icon")
			ui.SidebarItem(c, "/avatar", icons.User, "Avatar")
			ui.SidebarItem(c, "/video", icons.Video, "Video")
			ui.SidebarItem(c, "/audio", icons.Music, "Audio")
		})
		ui.SidebarSection(c, "Plugins", nil, func() {
			ui.SidebarItem(c, "/terminal", icons.Terminal, "Terminal")
		})
		ui.SidebarSection(c, "Desktop APIs", nil, func() {
			ui.SidebarItem(c, "/dialog", icons.Folder, "Dialog")
			ui.SidebarItem(c, "/notification", icons.Bell, "Notification")
			ui.SidebarItem(c, "/clipboard", icons.Clipboard, "Clipboard")
			ui.SidebarItem(c, "/shell", icons.Terminal, "Shell")
			ui.SidebarItem(c, "/shortcut", icons.Keyboard, "GlobalShortcut")
		})
		ui.SidebarSection(c, "反馈", nil, func() {
			ui.SidebarItem(c, "/status", icons.Activity, "Status")
			ui.SidebarItem(c, "/message", icons.Message, "Message")
			ui.SidebarItem(c, "/modal", icons.Dialog, "Modal")
		})
	}).Width(200).OnChange(func() {
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
		case r.Match("/choices"):
			app.pageChoices(c)
		case r.Match("/actions"):
			app.pageActions(c)
		case r.Match("/dragdrop"):
			app.pageDragDrop(c)
		case r.Match("/list"):
			app.pageList(c)
		case r.Match("/table"):
			app.pageTable(c)
		case r.Match("/tree"):
			app.pageTree(c)
		case r.Match("/outline"):
			app.pageOutline(c)
		case r.Match("/gridview"):
			app.pageGridView(c)
		case r.Match("/charts"):
			app.pageCharts(c)
		case r.Match("/image"):
			app.pageImage(c)
		case r.Match("/icon"):
			app.pageIcon(c)
		case r.Match("/avatar"):
			app.pageAvatar(c)
		case r.Match("/video"):
			app.pageVideo(c)
		case r.Match("/audio"):
			app.pageAudio(c)
		case r.Match("/terminal"):
			app.pageTerminal(c)
		case r.Match("/code"):
			app.pageCode(c)
		case r.Match("/document"):
			app.pageDocument(c)
		case r.Match("/status"):
			app.pageStatus(c)
		case r.Match("/message"):
			app.pageMessage(c)
		case r.Match("/modal"):
			app.pageModal(c)
		case r.Match("/dialog"):
			app.pageDialog(c)
		case r.Match("/notification"):
			app.pageNotification(c)
		case r.Match("/clipboard"):
			app.pageClipboard(c)
		case r.Match("/shell"):
			app.pageShell(c)
		case r.Match("/shortcut"):
			app.pageShortcut(c)
		default:
			ui.Text(c, "Not found")
		}
	})
}
