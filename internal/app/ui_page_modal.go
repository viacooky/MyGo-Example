package app

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

type pageModalData struct {
	DialogOpen bool
	AlertOpen  bool
	Name       string
	LastAction string
}

func (app *app) pageModal(c *ui.Context) {
	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Modal")
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Modal").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("对话框打开: %t", app.pageModalData.DialogOpen))
					ui.Text(c, fmt.Sprintf("最近操作: %s", app.pageModalData.LastAction))
				})

				ui.Text(c, "ui.Modal")
				if ui.Button(c, "打开对话框").Clicked() {
					app.pageModalData.DialogOpen = true
				}
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "AlertDialog").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("确认框打开: %t", app.pageModalData.AlertOpen))
				})

				ui.Text(c, "ui.AlertDialog")
				if ui.Button(c, "删除…").Clicked() {
					app.pageModalData.AlertOpen = true
				}
			})
	})

	// 对话框是覆盖层，放在页面内容之后构建
	ui.Modal(c, &app.pageModalData.DialogOpen, func() {
		ui.Text(c, "重命名").FontSize(18).Bold()
		ui.Text(c, "输入新的名称。").TextColor(c.Theme().TextMuted)
		ui.TextInput(c, &app.pageModalData.Name).AutoFocus()
		ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
			if ui.Button(c, "取消").Clicked() {
				app.pageModalData.DialogOpen = false
			}
			if ui.PrimaryButton(c, "确定").Clicked() {
				app.pageModalData.LastAction = fmt.Sprintf("重命名为 %q", app.pageModalData.Name)
				app.pageModalData.DialogOpen = false
			}
		})
	})

	switch ui.AlertDialog(c, &app.pageModalData.AlertOpen, "删除这条记录？", "删除后无法恢复。", "取消", "删除") {
	case ui.AlertDismissed:
		app.pageModalData.LastAction = "已取消删除"
	case 0:
		app.pageModalData.LastAction = "已取消删除"
	case 1:
		app.pageModalData.LastAction = "已删除"
	}
}
