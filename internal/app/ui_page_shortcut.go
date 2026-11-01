package app

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

// defaultShortcut 是页面首次打开时自动注册的快捷键
const defaultShortcut = "Ctrl+F1"

type pageShortcutData struct {
	Input       string // 录制得到的快捷键组合
	Current     string // 当前已注册的快捷键，空表示未注册
	Result      string
	Recording   bool // 正在等待按键
	PressCount  int
	initialized bool
}

func (app *app) pageShortcut(c *ui.Context) {
	d := &app.pageShortcutData
	if !d.initialized {
		d.initialized = true
		d.Input = defaultShortcut
		app.applyShortcut()
	}
	registered := d.Current != "" && mygo.GlobalShortcut.IsRegistered(d.Current)

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "GlobalShortcut")
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "注册与注销").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("已注册: %t", registered))
					ui.Text(c, fmt.Sprintf("当前快捷键: %s", orNone(d.Current)))
					ui.Text(c, fmt.Sprintf("触发次数: %d", d.PressCount))
					ui.Text(c, fmt.Sprintf("结果: %s", orNone(d.Result)))
				})

				ui.Text(c, "mygo.GlobalShortcut.Register")
				recorder := ui.Button(c, shortcutRecorderLabel(d)).Width(260).Focusable()
				if recorder.Clicked() {
					d.Recording = !d.Recording
					if d.Recording {
						recorder.Focus()
					}
				}
				recorder.HandleInput(func(ev ui.InputEvent) bool {
					return app.recordShortcut(ev)
				})
				ui.Row(c).Gap(8).Children(func() {
					if ui.Button(c, "应用录制的组合").Disabled(d.Input == "").Clicked() {
						app.applyShortcut()
					}
					if ui.Button(c, "注销").Disabled(!registered).Clicked() {
						app.unregisterShortcut()
					}
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "说明").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "点击录制按钮后直接按下组合键，例如 Ctrl+F1、Ctrl+Shift+K、Alt+Space。")
					ui.Text(c, "录制过程中按 Esc 取消；不按修饰键（功能键除外）的组合不会被接受。")
					ui.Text(c, "录制成功后会先注销旧快捷键再注册新的；注册成功后，无论在应用内还是切到其他窗口都生效。")
				})

				ui.Text(c, "Window.Show / Hide")
				if ui.Button(c, "立即切换窗口").JustifySelf(ui.Start).Clicked() {
					app.shortcutToggle()
				}
			})
	})
}

// shortcutRecorderLabel 返回录制按钮上显示的文案
func shortcutRecorderLabel(d *pageShortcutData) string {
	switch {
	case d.Recording:
		return "请按下快捷键…（Esc 取消）"
	case d.Input == "":
		return "点击录制快捷键"
	default:
		return d.Input
	}
}

// recordShortcut 处理录制按钮收到的按键事件，成功录制后立即注册
func (app *app) recordShortcut(ev ui.InputEvent) bool {
	d := &app.pageShortcutData
	if !d.Recording || ev.Kind != ui.InputKeyDown {
		return false
	}
	if ev.Repeat {
		return true
	}
	if ev.Key == ui.KeyEscape {
		d.Recording = false
		d.Result = "已取消录制"
		app.window.Update(func() {})
		return true
	}
	key, ok := acceleratorKeyName(ev.Key)
	if !ok {
		// 单独按下修饰键或无法识别的键，继续等待
		return true
	}
	if ev.Mods == 0 && !isFunctionKey(key) {
		d.Result = "请至少按住 Ctrl、Shift 或 Alt 再按一个键"
		app.window.Update(func() {})
		return true
	}
	d.Input = shortcutAccelerator(ev.Mods, key)
	d.Recording = false
	app.applyShortcut()
	app.window.Update(func() {})
	return true
}

// shortcutAccelerator 把修饰键与按键名组合成 mygo 接受的快捷键字符串
func shortcutAccelerator(mods ui.Modifiers, key string) string {
	parts := make([]string, 0, 5)
	if mods&ui.Ctrl != 0 {
		parts = append(parts, "Ctrl")
	}
	if mods&ui.Shift != 0 {
		parts = append(parts, "Shift")
	}
	if mods&ui.Alt != 0 {
		parts = append(parts, "Alt")
	}
	if mods&ui.Super != 0 {
		parts = append(parts, "Super")
	}
	return strings.Join(append(parts, key), "+")
}

// acceleratorKeyName 把 ui.Key 转为 mygo 快捷键里的按键名
func acceleratorKeyName(k ui.Key) (string, bool) {
	switch {
	case k >= ui.KeyA && k <= ui.KeyZ:
		return string(rune('A' + (k - ui.KeyA))), true
	case k >= ui.Key0 && k <= ui.Key9:
		return string(rune('0' + (k - ui.Key0))), true
	case k >= ui.KeyF1 && k <= ui.KeyF12:
		return fmt.Sprintf("F%d", int(k-ui.KeyF1)+1), true
	}
	name, ok := namedAcceleratorKeys[k]
	return name, ok
}

// isFunctionKey 判断按键名是否是功能键，这类键可以不搭配修饰键
func isFunctionKey(name string) bool {
	return len(name) > 1 && name[0] == 'F'
}

var namedAcceleratorKeys = map[ui.Key]string{
	ui.KeySpace:        "Space",
	ui.KeyEnter:        "Enter",
	ui.KeyTab:          "Tab",
	ui.KeyBackspace:    "Backspace",
	ui.KeyDelete:       "Delete",
	ui.KeyInsert:       "Insert",
	ui.KeyHome:         "Home",
	ui.KeyEnd:          "End",
	ui.KeyPageUp:       "PageUp",
	ui.KeyPageDown:     "PageDown",
	ui.KeyUp:           "Up",
	ui.KeyDown:         "Down",
	ui.KeyLeft:         "Left",
	ui.KeyRight:        "Right",
	ui.KeyMinus:        "-",
	ui.KeyEqual:        "=",
	ui.KeyComma:        ",",
	ui.KeyPeriod:       ".",
	ui.KeySlash:        "/",
	ui.KeySemicolon:    ";",
	ui.KeyQuote:        "'",
	ui.KeyBracketLeft:  "[",
	ui.KeyBracketRight: "]",
	ui.KeyBackslash:    `\`,
	ui.KeyBackquote:    "`",
}

// applyShortcut 注册录制得到的快捷键；已注册的旧快捷键会先被注销
func (app *app) applyShortcut() {
	d := &app.pageShortcutData
	acc := strings.TrimSpace(d.Input)
	if acc == "" {
		d.Result = "请先录制快捷键"
		return
	}
	if d.Current == acc && mygo.GlobalShortcut.IsRegistered(acc) {
		d.Result = "已注册 " + acc
		return
	}
	if d.Current != "" {
		mygo.GlobalShortcut.Unregister(d.Current)
		d.Current = ""
	}
	if err := mygo.GlobalShortcut.Register(acc, app.shortcutToggle); err != nil {
		d.Result = fmt.Sprintf("注册 %s 失败: %v", acc, err)
		return
	}
	d.Current = acc
	d.Result = "已注册 " + acc
}

// unregisterShortcut 注销当前快捷键
func (app *app) unregisterShortcut() {
	d := &app.pageShortcutData
	if d.Current != "" {
		mygo.GlobalShortcut.Unregister(d.Current)
	}
	d.Current = ""
	d.Result = "已注销快捷键"
}

// shortcutToggle 切换主窗口的显示状态；快捷键回调运行在主线程
func (app *app) shortcutToggle() {
	d := &app.pageShortcutData
	d.PressCount++
	app.toggleWindow()
	if app.window.IsVisible() {
		d.Result = "窗口已显示"
		return
	}
	if d.Current != "" {
		d.Result = fmt.Sprintf("窗口已隐藏，按 %s 恢复显示", d.Current)
		return
	}
	d.Result = "窗口已隐藏"
}
