package app

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

type pageCodeData struct {
	Editor *editor.State
	Diff   *editor.DiffState

	theme  *editor.Theme // 跟随窗口明暗缓存的主题
	dark   bool
	themed bool
}

const codeSample = `package main

import "fmt"

type Greeter struct {
	Name string
}

func (g Greeter) Greet() string {
	return "Hello, " + g.Name
}

func main() {
	items := []string{"go", "python", "rust"}
	for i, item := range items {
		fmt.Printf("%d: %s\n", i, item)
	}
	g := Greeter{Name: "MyGo"}
	fmt.Println(g.Greet())
}
`

const diffOldSample = `package main

func Add(a, b int) int {
	return a - b
}

func main() {
	println(Add(1, 2))
}
`

const diffNewSample = `package main

func Add(a, b int) int {
	return a + b
}

func Multiply(a, b int) int {
	return a * b
}

func main() {
	println(Add(1, 2))
	println(Multiply(3, 4))
}
`

func (app *app) pageCode(c *ui.Context) {
	ui.Row(c).Children(func() {
		ui.Text(c, "Code").FontSize(PageTitleFontSize).Margin(PageTitleBarMargin).Bold()
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Editor").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "可编辑的 Go 源码，支持语法高亮、多光标、折叠与查找")
				})

				ui.Text(c, "editor.View")
				if app.pageCodeData.Editor == nil {
					buf := editor.NewBuffer(codeSample)
					buf.SetLanguage(editor.LanguageFor("main.go"))
					app.pageCodeData.Editor = &editor.State{Buffer: buf}
				}
				app.pageCodeData.Editor.Theme = app.codeTheme(c)
				editor.View(c, app.pageCodeData.Editor).Height(260).FillWidth()
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "DiffView").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "新旧文本的并排对比，新的一侧可编辑")
				})

				ui.Text(c, "editor.DiffView")
				if app.pageCodeData.Diff == nil {
					lang := editor.LanguageFor("main.go")
					oldBuf := editor.NewBuffer(diffOldSample)
					oldBuf.SetLanguage(lang)
					newBuf := editor.NewBuffer(diffNewSample)
					newBuf.SetLanguage(lang)
					app.pageCodeData.Diff = &editor.DiffState{Old: oldBuf, New: newBuf, Editable: true}
				}
				app.pageCodeData.Diff.Theme = app.codeTheme(c)
				editor.DiffView(c, app.pageCodeData.Diff).Height(280).FillWidth()
			})
	})
}

// codeTheme 返回编辑器主题：配色沿用窗口主题，背景换成卡片色，与页面区分
func (app *app) codeTheme(c *ui.Context) *editor.Theme {
	theme := c.Theme()
	if app.pageCodeData.themed && app.pageCodeData.dark == theme.Dark {
		return app.pageCodeData.theme
	}
	th := editor.DefaultTheme(theme)
	th.Background = theme.Surface
	app.pageCodeData.theme = th
	app.pageCodeData.dark = theme.Dark
	app.pageCodeData.themed = true
	return th
}
