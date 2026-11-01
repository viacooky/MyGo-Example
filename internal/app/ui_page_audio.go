package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/audio"
	_ "github.com/egoist/mygo/audio/mp3"
	"github.com/egoist/mygo/ui"
)

// audioRefresh 是播放期间刷新进度的间隔
const audioRefresh = 200 * time.Millisecond

// audioTrack 是一条可供播放的示例音频
type audioTrack struct {
	Format string // 例：MP3
	Name   string
	Path   string
	Size   int64
}

type pageAudioData struct {
	Tracks []audioTrack
	Note   string

	Active   int // 当前曲目，-1 表示没有
	Playing  bool
	Position time.Duration
	Duration time.Duration
	Volume   float64
	Loop     bool
	Message  string

	initialized bool
	player      atomic.Pointer[audio.Player]
	ticking     atomic.Bool
}

func (app *app) pageAudio(c *ui.Context) {
	d := &app.pageAudioData
	if !d.initialized {
		d.initialized = true
		app.loadAudio()
	}
	t := c.Theme()

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Audio")
	})

	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Text(c, "audio 包内置 WAV 解码与流式播放；导入 audio/mp3 后同样支持 MP3。").
			TextColor(t.TextMuted).Margin(PageTitleBarMargin)
		if d.Note != "" {
			ui.Text(c, d.Note).TextColor(t.TextMuted).Margin(PageTitleBarMargin)
		}

		for i := range d.Tracks {
			if i > 0 {
				ui.Divider(c)
			}
			app.audioTrackRow(c, i)
		}

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "进度与定位").FontSize(16).Bold()
				ui.Column(c).Gap(8).Children(func() {
					if d.Active < 0 {
						ui.Text(c, "先播放一条音轨，再用下面的控件控制它。").
							TextColor(t.TextMuted)
						return
					}
					ui.Progress(c, app.audioFraction()).Label("播放进度")
					ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
						ui.Text(c, fmt.Sprintf("%s / %s",
							formatAudioTime(d.Position), formatAudioTime(d.Duration))).
							FontSize(12).TextColor(t.TextMuted)
						if ui.Button(c, "-10s").Clicked() {
							app.audioSeek(-10 * time.Second)
						}
						if ui.Button(c, "+10s").Clicked() {
							app.audioSeek(10 * time.Second)
						}
						if ui.Button(c, "停止").Clicked() {
							app.audioStop()
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
				ui.Text(c, "音量与循环").FontSize(16).Bold()
				ui.Column(c).Gap(8).Children(func() {
					if slider := ui.Slider(c, &d.Volume, 0, 1).Label("音量"); slider.Changed() {
						if p := d.player.Load(); p != nil {
							p.SetVolume(d.Volume)
						}
					}
					if ui.Switch(c, &d.Loop).Label("循环播放").Changed() {
						if p := d.player.Load(); p != nil {
							p.SetLoop(d.Loop)
						}
					}
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "说明").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "audio.OpenStream 返回暂停状态的 Player，边解码边播放，适合较长音频。")
					ui.Text(c, "Player 的 Play、Pause、Stop、Seek、SetVolume 与 SetLoop 可在任意 goroutine 调用；关闭 Player 会同时关闭打开的文件。")
					ui.Text(c, "MP3 解码来自 github.com/hajimehoshi/go-mp3，以空导入方式注册到 audio 包。")
					if d.Message != "" {
						ui.Text(c, d.Message).TextColor(ui.Hex("#ef4444"))
					}
				})
			})
	})
}

// audioTrackRow 展示一条音轨及其播放按钮
func (app *app) audioTrackRow(c *ui.Context, i int) {
	d := &app.pageAudioData
	t := c.Theme()
	track := d.Tracks[i]

	ui.Grid(c).Margin(PageTitleBarMargin).
		GapX(16).GapY(8).
		ColumnTracks(ui.FitContent(), ui.Fr(1)).
		AlignItems(ui.Baseline).
		Children(func() {
			ui.Text(c, track.Format).FontSize(16).Bold()
			ui.Column(c).Gap(8).Children(func() {
				ui.Text(c, fmt.Sprintf("%s · %s", track.Name, formatTableSize(track.Size))).
					FontSize(12).TextColor(t.TextMuted)

				ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
					label := "播放"
					if d.Active == i {
						label = "暂停"
						if !d.Playing {
							label = "继续"
						}
					}
					if ui.Button(c, label).Clicked() {
						if d.Active == i && d.Playing {
							app.audioPause()
						} else {
							app.audioPlay(i)
						}
					}
					stop := ui.Button(c, "停止")
					if d.Active != i {
						stop.Disabled(true)
					}
					if stop.Clicked() {
						app.audioStop()
					}
					if d.Active == i {
						ui.Progress(c, app.audioFraction()).Grow(1).Label("播放进度")
					} else {
						ui.Text(c, "未播放").FontSize(12).TextColor(t.TextMuted)
					}
				})
			})
		})
}

// loadAudio 记录示例音频的位置与大小，并在退出时关闭播放器
func (app *app) loadAudio() {
	d := &app.pageAudioData
	d.Volume = 0.8
	d.Active = -1

	dir, err := mygo.App.Path(mygo.PathResources)
	if err != nil {
		d.Note = fmt.Sprintf("无法定位资源目录: %v", err)
		return
	}
	for _, item := range []struct {
		format string
		name   string
	}{
		{"MP3", "audio-30s.mp3"},
		{"WAV", "audio-9s.wav"},
	} {
		path := filepath.Join(dir, item.name)
		info, err := os.Stat(path)
		if err != nil {
			d.Note = fmt.Sprintf("缺少示例音频 %s: %v", item.name, err)
			continue
		}
		d.Tracks = append(d.Tracks, audioTrack{
			Format: item.format,
			Name:   item.name,
			Path:   path,
			Size:   info.Size(),
		})
	}
	mygo.App.OnWillQuit(func(*mygo.QuitEvent) { app.audioClose() })
}

// audioPlay 播放第 i 条音轨，切换音轨时先关闭上一个播放器
func (app *app) audioPlay(i int) {
	d := &app.pageAudioData
	if i < 0 || i >= len(d.Tracks) {
		return
	}
	if d.Active != i {
		app.audioClose()

		file, err := os.Open(d.Tracks[i].Path)
		if err != nil {
			d.Message = fmt.Sprintf("无法打开 %s: %v", d.Tracks[i].Name, err)
			return
		}
		player, err := audio.OpenStream(file)
		if err != nil {
			file.Close()
			d.Message = fmt.Sprintf("无法解码 %s: %v", d.Tracks[i].Name, err)
			return
		}
		player.SetVolume(d.Volume).SetLoop(d.Loop)

		d.Active = i
		d.Position = 0
		d.Duration = player.Duration()
		d.Message = ""
		d.player.Store(player)
	}

	player := d.player.Load()
	if player == nil {
		return
	}
	player.Play()
	d.Playing = player.Playing()
	app.audioTick()
}

// audioPause 暂停当前音轨并记录播放位置
func (app *app) audioPause() {
	d := &app.pageAudioData
	player := d.player.Load()
	if player == nil {
		return
	}
	player.Pause()
	d.Position = player.Position()
	d.Playing = player.Playing()
}

// audioStop 停止当前音轨并回到开头
func (app *app) audioStop() {
	d := &app.pageAudioData
	player := d.player.Load()
	if player == nil {
		return
	}
	player.Stop()
	d.Position = 0
	d.Playing = player.Playing()
}

// audioSeek 把当前音轨移动到相对当前位置 delta 处
func (app *app) audioSeek(delta time.Duration) {
	d := &app.pageAudioData
	player := d.player.Load()
	if player == nil {
		return
	}
	target := d.Position + delta
	if d.Duration > 0 {
		target = min(max(target, 0), d.Duration)
	} else {
		target = max(target, 0)
	}
	player.Seek(target)
	d.Position = target
	d.Playing = player.Playing()
}

// audioClose 关闭当前播放器；Player 关闭时会一并关闭打开的文件
func (app *app) audioClose() {
	d := &app.pageAudioData
	if player := d.player.Swap(nil); player != nil {
		player.Close()
	}
	d.Active = -1
	d.Playing = false
	d.Position = 0
	d.Duration = 0
}

// audioTick 在播放期间定时刷新进度；ticker 随播放开始启动，暂停后退出
func (app *app) audioTick() {
	d := &app.pageAudioData
	if !d.ticking.CompareAndSwap(false, true) {
		return
	}
	go func() {
		ticker := time.NewTicker(audioRefresh)
		defer ticker.Stop()
		for range ticker.C {
			player := d.player.Load()
			if player == nil || !player.Playing() {
				d.ticking.Store(false)
				if app.window != nil {
					app.window.Update(func() {
						// 播放自然结束后同步状态，按钮才会回到「继续」
						d.Playing = false
						if p := d.player.Load(); p != nil {
							d.Playing = p.Playing()
							d.Position = p.Position()
							if duration := p.Duration(); duration > 0 {
								d.Duration = duration
							}
						}
					})
				}
				// 停止检查与恢复播放交错时，重新启动刷新，避免进度停在原处
				if next := d.player.Load(); next != nil && next.Playing() {
					app.audioTick()
				}
				return
			}
			position, duration := player.Position(), player.Duration()
			if app.window != nil {
				app.window.Update(func() {
					d.Position = position
					if duration > 0 {
						d.Duration = duration
					}
					d.Playing = true
				})
			}
		}
	}()
}

// audioFraction 返回当前音轨的播放进度，0 到 1；未知时长时为 -1
func (app *app) audioFraction() float64 {
	d := &app.pageAudioData
	if d.Duration <= 0 {
		return -1
	}
	return min(max(float64(d.Position)/float64(d.Duration), 0), 1)
}

// formatAudioTime 把时长格式化为 分:秒
func formatAudioTime(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	seconds := int(d.Seconds())
	return fmt.Sprintf("%d:%02d", seconds/60, seconds%60)
}
