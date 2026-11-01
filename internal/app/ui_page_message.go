package app

import (
	"fmt"

	"my-app/internal/icons"

	"github.com/egoist/mygo/ui"
)

type pageMessageData struct {
	Undone int // 通过 Toast 撤销的次数
}

func (app *app) pageMessage(c *ui.Context) {
	ui.Row(c).Children(func() {
		ui.Text(c, "Message").FontSize(PageTitleFontSize).Margin(PageTitleBarMargin).Bold()
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Tooltip").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "指针悬停或拿到键盘焦点时出现的提示")
				})

				ui.Text(c, "Element.Tooltip")
				ui.Row(c).Gap(8).Children(func() {
					ui.Button(c, "悬停查看提示").Tooltip("这是按钮的工具提示")
					ui.Icon(c, icons.Alert).Size(16, 16).Tooltip("图标也可以带工具提示")
					ui.Text(c, "文本提示").Tooltip("文本同样支持 Tooltip").TextColor(c.Theme().TextMuted)
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Toast").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "窗口底部短暂出现的提示，默认 4 秒后消失")
					ui.Text(c, fmt.Sprintf("已通过 Toast 撤销 %d 次", app.pageMessageData.Undone))
				})

				ui.Text(c, "Context.Toast")
				ui.Column(c).Gap(8).Children(func() {
					ui.Row(c).Gap(8).Children(func() {
						if ui.Button(c, "显示提示").Clicked() {
							c.Toast("操作成功")
						}
						if ui.Button(c, "带描述").Clicked() {
							c.AddToast(ui.Toast{Title: "上传失败", Description: "服务器无法连接，请稍后重试。"})
						}
						if ui.Button(c, "带操作").Clicked() {
							c.ToastAction("笔记已删除", "撤销", func() {
								app.pageMessageData.Undone++
							})
						}
					})
					ui.Row(c).Gap(8).Children(func() {
						if ui.Button(c, "常驻提示").Clicked() {
							c.AddToast(ui.Toast{Title: "正在同步…", Timeout: -1})
						}
						if ui.Button(c, "关闭全部提示").Clicked() {
							c.CloseToast("")
						}
					})
				})
			})
	})
}
