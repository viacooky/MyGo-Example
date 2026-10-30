package app

import "github.com/egoist/mygo/ui"

func (app *app) pageText(c *ui.Context) {
	ui.Row(c).Children(func() {
		ui.Text(c, "文本").FontSize(PageTitleFontSize).Margin(PageTitleBarMargin).Bold()
	})

	ui.Column(c).Selectable().Gap(12).Children(func() {
		ui.Text(c, "First paragraph.")
		ui.Text(c, "Second paragraph.")
	})

	ui.Column(c).Children(func() {
		ui.Text(c, "富文本").Bold()
		ui.RichText(c,
			ui.Span{Text: "Saved "},
			ui.Span{Text: "report.pdf", Weight: 600},
			ui.Span{Text: " to "},
			ui.Span{Text: "Documents", Color: c.Theme().Accent, Underline: true},
		)
		ui.Link(c, "首页", "/home")
	})

}
