package app

import (
	"github.com/egoist/mygo/ui"
)

func (app *app) pageText(c *ui.Context) {
	ui.Row(c).Children(func() {
		ui.Text(c, "Text").FontSize(PageTitleFontSize).Margin(PageTitleBarMargin).Bold()
	})

	ui.Grid(c).Margin(PageTitleBarMargin).
		ColumnTracks(ui.FitContent()).GapX(16).GapY(8).
		Children(func() {
			ui.Text(c, "这里是普通文本")
			ui.Text(c, "这里是加粗文本").Bold()
			ui.Text(c, "这里是24号文本").FontSize(24)
			ui.Text(c, "这里是斜体文本").Italic()
			ui.Text(c, "这里是下划线文本").Underline()
			ui.Text(c, "这里是删除线文本").Strikethrough()
			ui.Text(c, "这里是带颜色文本").TextColor(c.Theme().AccentPressed)
			ui.Text(c, "这里是边框的文本").Border(1, c.Theme().Border)
			ui.Text(c, "这里是圆角边框的文本").Border(1, c.Theme().Border).Radius(10)
			ui.Text(c, "这里是圆角边框的文本").Background(c.Theme().Surface)
			ui.Link(c, "这里是链接文本（百度）", "https://www.baidu.com/")

		})

}
