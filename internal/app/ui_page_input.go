package app

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
)

type pageInputData struct {
	initialized bool

	// 文本类控件全部绑定同一个字段，任意一处修改都会同步到其他控件
	Text string

	Number    float64   // 数字输入
	StepValue float64   // 步进滑块
	Date      time.Time // 日期与时间共用
	Color     ui.Color  // 色板与取色器共用
}

func (app *app) pageInput(c *ui.Context) {
	d := &app.pageInputData
	if !d.initialized {
		d.Number = 42
		d.StepValue = 5
		d.Date = time.Now()
		d.Color = ui.RGB(255, 255, 255)
		d.initialized = true
	}

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Input")
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		app.inputTextSection(c)
		ui.Divider(c)
		app.inputNumberSection(c)
		ui.Divider(c)
		app.inputDateSection(c)
		ui.Divider(c)
		app.inputColorSection(c)
	})
}

// inputTextSection 展示文本类控件：单行、多行、密码、搜索、自动补全与可编辑文本
func (app *app) inputTextSection(c *ui.Context) {
	ui.Grid(c).Margin(PageTitleBarMargin).
		GapX(16).GapY(8).
		ColumnTracks(ui.FitContent(), ui.Fr(1)).
		AlignItems(ui.Baseline).
		Children(func() {
			ui.Text(c, "Text").FontSize(16).Bold()
			ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
				ui.Text(c, fmt.Sprintf("输入的文本: %s", orNone(app.pageInputData.Text)))
				ui.Text(c, fmt.Sprintf("字数统计: %d", utf8.RuneCountInString(app.pageInputData.Text)))
			})

			ui.Text(c, "TextInput")
			ui.TextInput(c, &app.pageInputData.Text).Label("文本输入").AutoFocus()

			ui.Text(c, "With Placeholder")
			ui.TextInput(c, &app.pageInputData.Text).Label("带占位符").Placeholder("请输入")

			ui.Text(c, "Password")
			ui.TextInput(c, &app.pageInputData.Text).Label("密码").Password()

			ui.Text(c, "ReadOnly")
			ui.TextInput(c, &app.pageInputData.Text).Label("只读文本").ReadOnly(true)

			ui.Text(c, "Disabled")
			ui.TextInput(c, &app.pageInputData.Text).Label("禁用文本").Disabled(true)

			ui.Text(c, "TextArea")
			ui.TextArea(c, &app.pageInputData.Text).Label("多行文本").Height(160)

			ui.Text(c, "SearchField")
			ui.SearchField(c, &app.pageInputData.Text).Label("搜索")

			ui.Text(c, "Autocomplete")
			items := []string{"Amsterdam", "Berlin", "Paris", "Rome", "Madrid", "Tokyo"}
			ui.Autocomplete(c, &app.pageInputData.Text, items).Label("自动补全")

			ui.Text(c, "EditableText")
			ui.EditableText(c, &app.pageInputData.Text).Label("可编辑文本")
		})
}

// inputNumberSection 展示数字输入与步进滑块
func (app *app) inputNumberSection(c *ui.Context) {
	ui.Grid(c).Margin(PageTitleBarMargin).
		GapX(16).GapY(8).
		ColumnTracks(ui.FitContent(), ui.Fr(1)).
		AlignItems(ui.Baseline).
		Children(func() {
			ui.Text(c, "Number").FontSize(16).Bold()
			ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
				ui.Text(c, fmt.Sprintf("数字: %g", app.pageInputData.Number))
				ui.Text(c, fmt.Sprintf("步进值: %g", app.pageInputData.StepValue))
			})

			ui.Text(c, "NumberInput")
			ui.NumberInput(c, &app.pageInputData.Number, 1, 99, 0.1).Label("数字")

			ui.Text(c, "StepSlider")
			ui.StepSlider(c, &app.pageInputData.StepValue, 0, 10, 1).Label("步进滑块")
		})
}

// inputDateSection 展示日期与时间输入，两者绑定同一字段
func (app *app) inputDateSection(c *ui.Context) {
	ui.Grid(c).Margin(PageTitleBarMargin).
		GapX(16).GapY(8).
		ColumnTracks(ui.FitContent(), ui.Fr(1)).
		AlignItems(ui.Baseline).
		Children(func() {
			ui.Text(c, "Date & Time").FontSize(16).Bold()
			ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
				ui.Text(c, fmt.Sprintf("日期: %s", app.pageInputData.Date.Format("2006-01-02")))
				ui.Text(c, fmt.Sprintf("时间: %s", app.pageInputData.Date.Format("15:04:05")))
			})

			ui.Text(c, "DateInput")
			ui.DateInput(c, &app.pageInputData.Date).Label("日期")

			ui.Text(c, "TimeInput")
			ui.TimeInput(c, &app.pageInputData.Date).Label("时间")
		})
}

// inputColorSection 展示色板与取色器，两者绑定同一颜色
func (app *app) inputColorSection(c *ui.Context) {
	ui.Grid(c).Margin(PageTitleBarMargin).
		GapX(16).GapY(8).
		ColumnTracks(ui.FitContent(), ui.Fr(1)).
		AlignItems(ui.Baseline).
		Children(func() {
			ui.Text(c, "Color").FontSize(16).Bold()
			ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
				ui.Text(c, fmt.Sprintf("颜色: %+v", app.pageInputData.Color))
			})

			ui.Text(c, "ColorWell")
			ui.ColorWell(c, &app.pageInputData.Color).Label("颜色")

			ui.Text(c, "ColorPicker")
			ui.ColorPicker(c, &app.pageInputData.Color).Label("颜色")
		})
}
