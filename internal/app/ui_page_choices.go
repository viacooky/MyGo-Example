package app

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

type pageChoicesData struct {
	initialized bool

	ComboboxField string
	TokensField   []string

	Mail     bool
	Calendar bool
	Wifi     bool
	Size     string
	Volume   float64
	MinPrice float64
	MaxPrice float64
	Copies   float64
	Stars    int
}

func (app *app) pageChoices(c *ui.Context) {
	d := &app.pageChoicesData
	if !d.initialized {
		d.Size = "中"
		d.Volume = 50
		d.MinPrice = 200
		d.MaxPrice = 800
		d.Copies = 1
		d.Stars = 3
		d.initialized = true
	}
	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Choices")
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Combobox").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("选择的值: %s", app.pageChoicesData.ComboboxField))
					ui.Text(c, fmt.Sprintf("选择的令牌值: %+v", app.pageChoicesData.TokensField))
				})

				items := []string{"go", "python", "dotnet", "java", "javascript", "rust"}
				ui.Text(c, "Combobox")
				ui.Combobox(c, &app.pageChoicesData.ComboboxField, items)

				ui.Text(c, "Select")
				ui.Select(c, &app.pageChoicesData.ComboboxField, items)

				ui.Text(c, "TokenField")
				ui.TokenField(c, &app.pageChoicesData.TokensField, items)
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Checkbox").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("邮件: %t, 日历: %t", app.pageChoicesData.Mail, app.pageChoicesData.Calendar))
				})

				ui.Text(c, "Checkbox")
				ui.Column(c).Gap(8).Children(func() {
					ui.Checkbox(c, &app.pageChoicesData.Mail, "接收邮件")
					ui.Checkbox(c, &app.pageChoicesData.Calendar, "同步日历")
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Switch").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("Wi-Fi: %t", app.pageChoicesData.Wifi))
				})

				ui.Text(c, "Switch")
				ui.Row(c).Gap(12).Children(func() {
					ui.Text(c, "Wi-Fi")
					ui.Switch(c, &app.pageChoicesData.Wifi).Label("Wi-Fi")
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Radio").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("尺寸: %s", app.pageChoicesData.Size))
				})

				ui.Text(c, "RadioGroup")
				ui.RadioGroup(c, func() {
					ui.Radio(c, &app.pageChoicesData.Size, "小", "小")
					ui.Radio(c, &app.pageChoicesData.Size, "中", "中")
					ui.Radio(c, &app.pageChoicesData.Size, "大", "大")
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Slider").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("音量: %.0f", app.pageChoicesData.Volume))
				})

				ui.Text(c, "Slider")
				ui.Slider(c, &app.pageChoicesData.Volume, 0, 100).Label("音量")
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "RangeSlider").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("价格: %.0f – %.0f", app.pageChoicesData.MinPrice, app.pageChoicesData.MaxPrice))
				})

				ui.Text(c, "RangeSlider")
				ui.RangeSlider(c, &app.pageChoicesData.MinPrice, &app.pageChoicesData.MaxPrice, 0, 1000, 10).Label("价格")
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Stepper").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("份数: %.0f", app.pageChoicesData.Copies))
				})

				ui.Text(c, "Stepper")
				ui.Stepper(c, &app.pageChoicesData.Copies, 1, 99, 1).Label("份数")
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Rating").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("评分: %d", app.pageChoicesData.Stars))
				})

				ui.Text(c, "Rating")
				ui.Rating(c, &app.pageChoicesData.Stars, 5).Label("评分")
			})
	})
}
