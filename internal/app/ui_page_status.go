package app

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/egoist/mygo/ui"
)

type pageStatusData struct {
	Progress    float64 // 下载进度，0 到 1
	Speed       float64 // 下载速度，MB/s
	Downloaded  float64 // 已下载，MB
	Downloading bool

	SpinnerDownloaded  float64 // spinner 示例的已下载量，MB
	SpinnerSpeed       float64 // spinner 示例的下载速度，MB/s
	SpinnerDownloading bool

	DiskUsed float64
	Unread   int
}

const downloadSize = 20.0 // 模拟文件大小，MB

func (app *app) pageStatus(c *ui.Context) {
	ui.Row(c).Children(func() {
		ui.Text(c, "Status").FontSize(PageTitleFontSize).Margin(PageTitleBarMargin).Bold()
	})
	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Progress").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("进度: %.0f%%", app.pageStatusData.Progress*100))
					ui.Text(c, fmt.Sprintf("速度: %.1f MB/s", app.pageStatusData.Speed))
					ui.Text(c, fmt.Sprintf("已下载: %.1f / %.0f MB", app.pageStatusData.Downloaded, downloadSize))
					ui.Text(c, fmt.Sprintf("下载中: %t", app.pageStatusData.Downloading))
				})

				ui.Text(c, "Progress")
				ui.Column(c).Gap(8).Children(func() {
					ui.Row(c).Gap(8).Children(func() {
						button := ui.Button(c, "开始下载")
						if app.pageStatusData.Downloading {
							button.Disabled(true)
						}
						if button.Clicked() {
							app.pageStatusData.Progress = 0
							app.pageStatusData.Speed = 0
							app.pageStatusData.Downloaded = 0
							app.pageStatusData.Downloading = true
							go app.simulateDownload(app.updateDownload)
						}
					})
					ui.Progress(c, app.pageStatusData.Progress).Label("下载进度")
				})

				ui.Text(c, "Unknown length")
				ui.Progress(c, -1).Label("连接中")
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Spinner").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "加载中的旋转指示器")
				})

				ui.Text(c, "Spinner")
				ui.Row(c).Gap(8).Children(func() {
					ui.Spinner(c).Label("加载中")
					ui.Text(c, "正在加载…").TextColor(c.Theme().TextMuted)
				})

				ui.Text(c, "Download Spinner")
				ui.Column(c).Gap(8).Children(func() {
					ui.Row(c).Gap(8).Children(func() {
						button := ui.Button(c, "开始下载")
						if app.pageStatusData.SpinnerDownloading {
							button.Disabled(true)
						}
						if button.Clicked() {
							app.pageStatusData.SpinnerDownloaded = 0
							app.pageStatusData.SpinnerSpeed = 0
							app.pageStatusData.SpinnerDownloading = true
							go app.simulateDownload(app.updateSpinnerDownload)
						}
					})
					ui.Row(c).Gap(8).Children(func() {
						switch {
						case app.pageStatusData.SpinnerDownloading:
							ui.Spinner(c).Label("下载中")
							ui.Text(c, fmt.Sprintf("正在下载 %.1f / %.0f MB · %.1f MB/s",
								app.pageStatusData.SpinnerDownloaded, downloadSize, app.pageStatusData.SpinnerSpeed)).
								TextColor(c.Theme().TextMuted)
						case app.pageStatusData.SpinnerDownloaded >= downloadSize:
							ui.Text(c, fmt.Sprintf("下载完成，共 %.0f MB", downloadSize)).
								TextColor(c.Theme().TextMuted)
						default:
							ui.Text(c, "点击「开始下载」模拟下载").TextColor(c.Theme().TextMuted)
						}
					})
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Meter").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("磁盘占用: %.0f%%", app.pageStatusData.DiskUsed))
				})

				ui.Text(c, "Meter")
				ui.Column(c).Gap(8).Children(func() {
					ui.Meter(c, app.pageStatusData.DiskUsed, 0, 100, &ui.MeterLevels{Warning: 70, Critical: 90}).Label("磁盘占用")
					ui.Slider(c, &app.pageStatusData.DiskUsed, 0, 100).Label("磁盘占用百分比")
					ui.NumberInput(c, &app.pageStatusData.DiskUsed, 0, 100, 1).Label("磁盘占用百分比")
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Badge").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "数量徽标与强调色徽标")
				})

				ui.Text(c, "Badge")
				ui.Row(c).Gap(8).Children(func() {
					ui.Text(c, "收件箱")
					ui.Badge(c, fmt.Sprint(app.pageStatusData.Unread))
					ui.Badge(c, "New").Background(c.Theme().Accent).TextColor(c.Theme().AccentText)
				})
			})
	})
}

// simulateDownload 在后台按随机速度推进下载，每 100ms 把进度、速度与已下载量
// 交回主线程，由 report 写入各自示例的状态。
func (app *app) simulateDownload(report func(speed, downloaded float64, done bool)) {
	downloaded := 0.0
	for downloaded < downloadSize {
		time.Sleep(100 * time.Millisecond)

		speed := 3 + 6*rand.Float64() // 3 到 9 MB/s
		downloaded = min(downloaded+speed*0.1, downloadSize)
		// 先在后台线程取好快照，交给主线程的闭包不再读循环变量
		mb := downloaded
		done := downloaded >= downloadSize

		app.window.Update(func() {
			report(speed, mb, done)
		})
	}
}

// updateDownload 把下载快照写入进度示例的状态
func (app *app) updateDownload(speed, downloaded float64, done bool) {
	app.pageStatusData.Progress = downloaded / downloadSize
	app.pageStatusData.Speed = speed
	app.pageStatusData.Downloaded = downloaded
	if done {
		app.pageStatusData.Downloading = false
		app.pageStatusData.Speed = 0
	}
}

// updateSpinnerDownload 把下载快照写入 spinner 示例的状态
func (app *app) updateSpinnerDownload(speed, downloaded float64, done bool) {
	app.pageStatusData.SpinnerDownloaded = downloaded
	app.pageStatusData.SpinnerSpeed = speed
	if done {
		app.pageStatusData.SpinnerDownloading = false
		app.pageStatusData.SpinnerSpeed = 0
	}
}
