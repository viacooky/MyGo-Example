package app

import (
	"fmt"

	"my-app/internal/icons"

	"github.com/egoist/mygo/ui"
)

type pageActionsData struct {
	ButtonCount   int
	SortBy        string
	LastAction    string
	Notifications bool
	AutoSave      bool
}

func (app *app) pageActions(c *ui.Context) {
	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Actions")
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Button").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("点击次数: %d", app.pageActionsData.ButtonCount))
				})

				ui.Text(c, "Button")
				if ui.Button(c, "点击我").JustifySelf(ui.Start).Clicked() {
					app.pageActionsData.ButtonCount++
				}

				ui.Text(c, "With Icon")
				ui.Button(c, "").Label("发送").JustifySelf(ui.Start).Children(func() { ui.Icon(c, icons.Send).Size(16, 16) })

				ui.Text(c, "Disabled")
				ui.Button(c, "禁用状态").JustifySelf(ui.Start).Disabled(true)
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "MenuButton").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("当前排序: %s", app.pageActionsData.SortBy))
				})

				if app.pageActionsData.SortBy == "" {
					app.pageActionsData.SortBy = "名称"
				}
				ui.Text(c, "MenuButton")
				ui.MenuButton(c, app.pageActionsData.SortBy, func(m *ui.Menu) {
					for _, by := range []string{"名称", "日期", "大小"} {
						if m.Item(by).Chosen() {
							app.pageActionsData.SortBy = by
						}
					}
				}).JustifySelf(ui.Start)

				ui.Text(c, "With Separator")
				ui.MenuButton(c, "更多操作", func(m *ui.Menu) {
					if m.Item("刷新").Chosen() {
						app.pageActionsData.ButtonCount = 0
					}
					m.Separator()
					if m.Item("清空").Disabled(true).Chosen() {
						app.pageActionsData.ButtonCount = 0
					}
				}).JustifySelf(ui.Start)
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "ContextMenu").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("最近操作: %s", app.pageActionsData.LastAction))
				})

				ui.Text(c, "右键区域")
				t := c.Theme()
				ui.Row(c).Padding(8, 16).Radius(6).Background(t.Surface).
					JustifySelf(ui.Start).
					ContextMenu(func(m *ui.Menu) {
						if m.Item("复制").Chosen() {
							app.pageActionsData.LastAction = "复制"
						}
						if m.Item("重命名").Chosen() {
							app.pageActionsData.LastAction = "重命名"
						}
						m.Separator()
						if m.Item("删除").Chosen() {
							app.pageActionsData.LastAction = "删除"
						}
					}).
					Children(func() { ui.Text(c, "在此区域右键打开菜单") })
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Toggle").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("通知: %t", app.pageActionsData.Notifications))
					ui.Text(c, fmt.Sprintf("自动保存: %t", app.pageActionsData.AutoSave))
				})

				ui.Text(c, "Toggle")
				styledToggle(c, &app.pageActionsData.Notifications, "开启通知")

				ui.Text(c, "Toggle")
				styledToggle(c, &app.pageActionsData.AutoSave, "自动保存")
			})
	})
}

// styledToggle 用 ToggleBase 自绘开关：选中时用强调色，未选中时用普通表面色
func styledToggle(c *ui.Context, on *bool, label string) {
	t := c.Theme()
	bg, fg := t.Surface, t.Text
	if *on {
		bg, fg = t.Accent, t.AccentText
	}
	toggle := ui.ToggleBase(c, on).Label(label).Padding(6, 12).Radius(6).Background(bg).JustifySelf(ui.Start)
	if !*on && toggle.Hovered() {
		toggle.Background(t.SurfaceHover)
	}
	toggle.Children(func() { ui.Text(c, label).TextColor(fg) })
}
