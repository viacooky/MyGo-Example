package app

import (
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type counter struct{ n int }

// app is the state the window shows. Its view builds the interface from
// it, on the main thread, whenever the window needs a frame: after input,
// after Window.Update, and while something animates.
type app struct {
	window *mygo.Window // 主窗口
	router *ui.Router   // 页面路由

	// 状态
	counter counter
	text    string
}

func Run() error {
	app := &app{router: ui.NewRouter("/home")} // 初始化路由
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
	})

	return mygo.App.Run()
}
