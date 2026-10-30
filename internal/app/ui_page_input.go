package app

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
)

type pageInputData struct {
	TextField     string
	EditableText  string
	NumberField   float64
	ComboboxField string
	TokensField   []string
	TimeField     time.Time
	ColorField    ui.Color
}

func (app *app) pageInput(c *ui.Context) {
	ui.Row(c).Children(func() {
		ui.Text(c, "Input").FontSize(PageTitleFontSize).Margin(PageTitleBarMargin).Bold()
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "TextInput").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).
					Children(func() {
						ui.Text(c, fmt.Sprintf("输入的文本: %s", app.pageInputData.TextField))
						ui.Text(c, fmt.Sprintf("字数统计: %d", utf8.RuneCountInString(app.pageInputData.TextField)))
					})

				ui.Text(c, "TextInput")
				ui.TextInput(c, &app.pageInputData.TextField).AutoFocus()

				ui.Text(c, "With Placeholder")
				ui.TextInput(c, &app.pageInputData.TextField).Placeholder("请输入")

				ui.Text(c, "TextArea")
				ui.TextArea(c, &app.pageInputData.TextField).Height(160)

				ui.Text(c, "Password")
				ui.TextInput(c, &app.pageInputData.TextField).Password()

				ui.Text(c, "ReadOnly")
				ui.TextInput(c, &app.pageInputData.TextField).ReadOnly(true)

				ui.Text(c, "Disabled")
				ui.TextInput(c, &app.pageInputData.TextField).Disabled(true)

				ui.Text(c, "Autocomplete")
				items := []string{"Amsterdam", "Berlin", "Paris", "Rome", "Madrid", "Tokyo"}
				ui.Autocomplete(c, &app.pageInputData.TextField, items)

				ui.Text(c, "EditableText")
				if app.pageInputData.EditableText == "" {
					app.pageInputData.EditableText = "双击进入编辑，回车或失焦出发值改变事件.pdf"
				}
				ui.EditableText(c, &app.pageInputData.EditableText)

			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "NumberInput").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).
					Children(func() {
						ui.Text(c, fmt.Sprintf("输入的数字: %f", app.pageInputData.NumberField))

					})

				ui.Text(c, "NumberInput")
				ui.NumberInput(c, &app.pageInputData.NumberField, 1, 99, 0.1)
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Combobox").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).
					Children(func() {
						ui.Text(c, fmt.Sprintf("选择的值: %s", app.pageInputData.ComboboxField))
						ui.Text(c, fmt.Sprintf("选择的令牌值: %+v", app.pageInputData.TokensField))

					})

				items := []string{"go", "python", "dotnet", "java", "javascript", "rust"}
				ui.Text(c, "Combobox")
				ui.Combobox(c, &app.pageInputData.ComboboxField, items)

				ui.Text(c, "Select")
				ui.Select(c, &app.pageInputData.ComboboxField, items)

				ui.Text(c, "TokenField")
				ui.TokenField(c, &app.pageInputData.TokensField, items)
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "DateInput").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).
					Children(func() {
						timeStr := app.pageInputData.TimeField.Format("2006-01-02 15:04:05")
						ui.Text(c, fmt.Sprintf("Time: %s", timeStr))
						ui.Text(c, "暂时不能修改控件中的日期/时间格式").TextColor(ui.Hex("#FF0000"))
					})

				if app.pageInputData.TimeField.IsZero() {
					app.pageInputData.TimeField = time.Now()
				}
				ui.Text(c, "DateInput")
				ui.DateInput(c, &app.pageInputData.TimeField)

				ui.Text(c, "TimeInput")
				ui.TimeInput(c, &app.pageInputData.TimeField)
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Color").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).
					Children(func() {
						ui.Text(c, fmt.Sprintf("Color: %+v", app.pageInputData.ColorField))
					})

				if app.pageInputData.ColorField == ui.Transparent {
					app.pageInputData.ColorField = ui.RGB(255, 255, 255)
				}
				ui.Text(c, "ColorWell")
				ui.ColorWell(c, &app.pageInputData.ColorField)

				ui.Text(c, "ColorPicker")
				ui.ColorPicker(c, &app.pageInputData.ColorField)

			})
	})
}
