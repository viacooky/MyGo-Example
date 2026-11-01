package app

import (
	"errors"
	"fmt"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

type pageNotificationData struct {
	Title     string
	Body      string
	Silent    bool
	Result    string
	Clicked   bool
	Supported *bool              // 平台是否支持桌面通知，首次构建时探测
	note      *mygo.Notification // 当前可点击的通知，用于关闭
}

func (app *app) pageNotification(c *ui.Context) {
	d := &app.pageNotificationData
	if d.Title == "" {
		d.Title = "导出完成"
		d.Body = "文件已保存到下载目录。"
	}
	if d.Supported == nil {
		supported := mygo.NotificationsSupported()
		d.Supported = &supported
	}

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Notification")
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "NewNotification").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("平台支持: %t", *d.Supported))
					ui.Text(c, fmt.Sprintf("结果: %s", orNone(d.Result)))
				})

				ui.Text(c, "消息标题")
				ui.TextInput(c, &d.Title).Label("标题")
				ui.Text(c, "消息正文")
				ui.TextInput(c, &d.Body).Label("正文")
				ui.Text(c, "静音")
				ui.Switch(c, &d.Silent).Label("静音")
				if ui.Button(c, "发送通知").JustifySelf(ui.Start).Clicked() {
					go app.showNotification(d.Title, d.Body, d.Silent)
				}
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "点击与关闭").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("通知被点击: %t", d.Clicked))
					ui.Text(c, "同一 ID 的通知会替换上一条")
				})

				ui.Text(c, "Notification.OnClick / Close")
				ui.Row(c).Gap(8).Children(func() {
					if ui.Button(c, "发送可点击通知").Clicked() {
						app.sendClickableNotification()
					}
					if ui.Button(c, "关闭通知").Disabled(d.note == nil).Clicked() {
						d.note.Close()
						d.note = nil
						d.Result = "已关闭通知"
					}
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "ClearNotifications").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "清除本应用的全部通知")
				})

				ui.Text(c, "mygo.ClearNotifications")
				if ui.Button(c, "清除全部").JustifySelf(ui.Start).Clicked() {
					mygo.ClearNotifications()
					d.Result = "已清除全部通知"
				}
			})
	})
}

// showNotification 在后台发送一条普通通知；标题与正文由调用方在主线程先取好
func (app *app) showNotification(title, body string, silent bool) {
	d := &app.pageNotificationData
	if !*d.Supported {
		app.window.Update(func() { d.Result = "当前平台不支持桌面通知" })
		return
	}
	n := mygo.NewNotification(mygo.NotificationOptions{
		Title:  title,
		Body:   body,
		Silent: silent,
	})
	err := n.Show()
	result := "已发送"
	switch {
	case errors.Is(err, mygo.ErrNotificationsDenied):
		result = "通知权限被拒绝"
	case err != nil:
		result = fmt.Sprintf("发送失败: %v", err)
	}
	app.window.Update(func() { d.Result = result })
}

// sendClickableNotification 创建一条带 ID 与点击回调的通知
func (app *app) sendClickableNotification() {
	d := &app.pageNotificationData
	n := mygo.NewNotification(mygo.NotificationOptions{
		ID:    "mygo-lean-click",
		Group: "mygo-lean",
		Title: "任务完成",
		Body:  "点击这条通知查看详情。",
	})
	n.OnClick(func() {
		app.window.Update(func() {
			d.Clicked = true
			d.Result = "通知被点击"
		})
	})
	d.note = n
	d.Result = "已发送，等待点击"
	go func() {
		if err := n.Show(); err != nil {
			app.window.Update(func() { d.Result = fmt.Sprintf("发送失败: %v", err) })
		}
	}()
}
