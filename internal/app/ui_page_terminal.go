package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/plugins/terminal"
	"github.com/egoist/mygo/ui"
)

type pageTerminalData struct {
	term        *terminal.Terminal
	err         error
	title       string
	exitCode    int
	exited      bool
	mu          sync.Mutex
	initialized bool
}

func (app *app) pageTerminal(c *ui.Context) {
	d := &app.pageTerminalData
	if !d.initialized {
		d.initialized = true
		app.restartTerminal()
	}

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Terminal")
	})

	ui.Column(c).Fill().Padding(0, 16, 32, 16).Gap(8).Children(func() {
		if d.err != nil {
			ui.Text(c, "无法启动终端").FontSize(16).Bold()
			ui.Text(c, d.err.Error()).TextColor(ui.Hex("#dc2626"))
			ui.Text(c, "终端插件依赖 libghostty-vt 原生库，项目 third_party/ghostty-vt 已附带开发用的副本，打包后使用随包资源，均无需联网下载。").
				TextColor(c.Theme().TextMuted)
			ui.Row(c).Children(func() {
				if ui.Button(c, "重试").Clicked() {
					app.restartTerminal()
				}
			})
			return
		}

		d.mu.Lock()
		title, code, exited := d.title, d.exitCode, d.exited
		d.mu.Unlock()
		cols, rows := d.term.Size()

		status := fmt.Sprintf("%d×%d", cols, rows)
		if title != "" {
			status = fmt.Sprintf("%s · %s", title, status)
		}
		if exited {
			status = fmt.Sprintf("%s · 已退出 (%d)", status, code)
		}
		ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
			ui.Text(c, status).TextColor(c.Theme().TextMuted).Grow(1)
			if ui.Button(c, "重启").Clicked() {
				app.restartTerminal()
			}
		})
		if exited {
			ui.Text(c, "程序已退出，点击「重启」开始新的会话。").TextColor(c.Theme().TextMuted)
		}
		terminal.View(c, d.term).Fill().Grow(1).AutoFocus().
			Background(c.Theme().Surface).
			Border(1, c.Theme().Border).
			Radius(8)
	})
}

// restartTerminal 结束当前终端并启动一个新的默认 Shell 会话
func (app *app) restartTerminal() {
	d := &app.pageTerminalData
	if d.term != nil {
		_ = d.term.Close()
		d.term = nil
	}
	d.mu.Lock()
	d.title, d.exitCode, d.exited = "", 0, false
	d.mu.Unlock()
	d.err = nil
	useLocalNative()

	term, err := terminal.New(terminal.Options{
		Cursor:      terminal.CursorBar,
		Transparent: true, // 背景留给视图，跟随页面的主题色
		OnTitle: func(title string) {
			d.mu.Lock()
			d.title = title
			d.mu.Unlock()
			app.window.Update(func() {})
		},
		OnExit: func(code int) {
			d.mu.Lock()
			d.exitCode = code
			d.exited = true
			d.mu.Unlock()
			app.window.Update(func() {})
		},
	})
	if err != nil {
		d.err = fmt.Errorf("创建终端失败: %w", err)
		return
	}
	d.term = term
}

// closeTerminal 在退出应用时挂断终端里的程序
func (app *app) closeTerminal() {
	d := &app.pageTerminalData
	if d.term != nil {
		_ = d.term.Close()
		d.term = nil
	}
}

// useLocalNative 让开发运行优先加载项目内 third_party/ghostty-vt 的原生库，
// 避免未打包时首次运行从网络下载；打包后的应用会命中随包资源，不受影响。
func useLocalNative() {
	if os.Getenv("MYGO_GHOSTTY_VT") != "" {
		return
	}
	name, ok := localNativeName()
	if !ok {
		return
	}
	var dirs []string
	if dir, err := mygo.App.Path(mygo.PathResources); err == nil {
		dirs = append(dirs, filepath.Join(filepath.Dir(dir), "third_party", "ghostty-vt"))
	}
	dirs = append(dirs, filepath.Join("third_party", "ghostty-vt"))
	for _, dir := range dirs {
		path, err := filepath.Abs(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			_ = os.Setenv("MYGO_GHOSTTY_VT", path)
			return
		}
	}
}

// localNativeName 返回当前平台对应的 libghostty-vt 文件名
func localNativeName() (string, bool) {
	switch runtime.GOOS {
	case "windows":
		return "ghostty-vt.dll", true
	case "darwin":
		return "libghostty-vt.dylib", true
	case "linux":
		return "libghostty-vt.so", true
	}
	return "", false
}
