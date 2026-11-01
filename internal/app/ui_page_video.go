package app

import (
	"github.com/egoist/mygo/ui"
)

func (app *app) pageVideo(c *ui.Context) {
	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Video")
	})

	ui.Column(c).Fill().Padding(0, 16, 64, 16).Gap(8).Children(func() {
		ui.Text(c, "暂时无法支持").FontSize(16).Bold()
		ui.Text(c, "MyGo v0.7.4 的 ui.Video 只接受自定义 VideoDecoder，内置仅实现 GIF。").
			TextColor(c.Theme().TextMuted)
		ui.Text(c, "H.264、HEVC、VP9、AV1 等系统编解码器（Media Foundation / AVFoundation / GStreamer）尚未提供，因此暂时无法内嵌播放 mp4/webm。").
			TextColor(c.Theme().TextMuted)
	})
}
