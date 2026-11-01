package app

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// doubleClickInterval 是判定托盘双击的最大间隔，与系统默认的双击时间一致
const doubleClickInterval = 500 * time.Millisecond

type counter struct{ n int }

// app is the state the window shows. Its view builds the interface from
// it, on the main thread, whenever the window needs a frame: after input,
// after Window.Update, and while something animates.
type app struct {
	window *mygo.Window // 主窗口
	router *ui.Router   // 页面路由

	breadcrumbChosen int // 标题面包屑当前选中的层级，页面标题统一使用

	lastTrayClick time.Time // 托盘上一次单击的时间，只在主线程访问

	// 页面数据，由各页面在首次构建时初始化
	pageInputData        pageInputData
	pageListData         pageListData
	pageActionsData      pageActionsData
	pageChoicesData      pageChoicesData
	pageCodeData         pageCodeData
	pageDocumentData     pageDocumentData
	pageDialogData       pageDialogData
	pageDragDropData     pageDragDropData
	pageGridViewData     pageGridViewData
	pageMessageData      pageMessageData
	pageModalData        pageModalData
	pageNotificationData pageNotificationData
	pageClipboardData    pageClipboardData
	pageShellData        pageShellData
	pageShortcutData     pageShortcutData
	pageOutlineData      pageOutlineData
	pageStatusData       pageStatusData
	pageTableData        pageTableData
	pageTreeData         pageTreeData

	pageAudioData    pageAudioData
	pageAvatarData   pageAvatarData
	pageChartsData   pageChartsData
	pageImageData    pageImageData
	pageTerminalData pageTerminalData
}

func Run() error {
	app := New()
	mygo.App.WhenReady(func() {
		app.window = mygo.NewWindow(mygo.WindowOptions{
			Title:         mygo.App.Name(),
			Width:         800,
			Height:        600,
			MinWidth:      800,
			MinHeight:     600,
			StateKey:      "main",
			TitleBarStyle: mygo.TitleBarHidden,
			Content:       ui.View(app.view),
		})
		app.setupTray()
		mygo.App.OnWillQuit(func(*mygo.QuitEvent) { app.closeTerminal() })
	})

	return mygo.App.Run()
}

func New() *app {
	return &app{
		router: ui.NewRouter("/home"),
	}
}

// setupTray 创建托盘图标：最小化时把窗口隐藏到托盘，可从托盘的菜单恢复或退出
func (app *app) setupTray() {
	dir, err := mygo.App.Path(mygo.PathResources)
	if err != nil {
		return
	}
	icon, err := os.ReadFile(filepath.Join(dir, "icon.png"))
	if err != nil {
		slog.Error(fmt.Sprintf("tray icon: %v", err))
		return
	}
	tray, err := mygo.NewTray(mygo.TrayOptions{
		ID:      mygo.App.Name(),
		Icon:    icon,
		ToolTip: mygo.App.Name(),
		Menu: mygo.NewMenu([]*mygo.MenuItem{
			{Label: "显示主窗口", Click: func(*mygo.MenuItem, *mygo.Window) { app.showWindow() }},
			mygo.Separator(),
			{Role: mygo.RoleQuit},
		}),
	})
	if err != nil {
		// 没有托盘（如缺少 libayatana-appindicator3 的 Linux）时保持普通的最小化行为
		slog.Error(fmt.Sprintf("tray: %v", err))
		return
	}
	mygo.App.OnWillQuit(func(*mygo.QuitEvent) { tray.Destroy() })
	app.window.OnMinimize(func() { app.window.Hide() })
	tray.OnClick(app.trayClicked)
}

// trayClicked 在系统双击间隔内收到第二次单击时，切换主窗口的显示
func (app *app) trayClicked() {
	now := time.Now()
	if now.Sub(app.lastTrayClick) <= doubleClickInterval {
		app.lastTrayClick = time.Time{}
		app.toggleWindow()
		return
	}
	app.lastTrayClick = now
}

// toggleWindow 显示时隐藏，隐藏时从托盘恢复
func (app *app) toggleWindow() {
	if app.window.IsVisible() {
		app.window.Hide()
		return
	}
	app.showWindow()
}

// showWindow 从托盘恢复并显示主窗口
func (app *app) showWindow() {
	app.window.Restore()
	app.window.Show()
}
