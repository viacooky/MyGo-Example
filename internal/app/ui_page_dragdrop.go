package app

import (
	"fmt"
	"os"
	"path/filepath"

	"my-app/internal/icons"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/transfer"
	"github.com/egoist/mygo/ui"
)

// maxDroppedFiles 是页面保留的文件记录上限，避免长时间拖放后无限增长
const maxDroppedFiles = 20

// droppedFile 是拖入的文件或文件夹，元数据在拖入时就读取好
type droppedFile struct {
	Path string
	Size int64 // 目录或读取失败时为 -1
	Dir  bool
}

// pageDragDropData 保存拖放页的接收记录与拖出用的示例数据
type pageDragDropData struct {
	Files       []droppedFile
	Dropped     int // 累计接收的文件数量
	Source      transfer.Data
	SourcePath  string
	SourceNote  string
	initialized bool
}

func (app *app) pageDragDrop(c *ui.Context) {
	d := &app.pageDragDropData
	if !d.initialized {
		d.initialized = true
		app.loadDragSource()
	}
	t := c.Theme()

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "DragDrop")
	})

	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "拖入文件").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("已接收 %d 个文件", d.Dropped))
					ui.Text(c, "从资源管理器或 Finder 把文件（可多个）拖到下面的区域。").
						TextColor(t.TextMuted)
				})

				ui.Text(c, "ui.Element.DroppedFiles \rFileDragOver")
				zone := ui.Box(c).Width(360).Height(160).Radius(8).
					Background(t.Surface).Border(1, t.Border).Center().JustifySelf(ui.Start)
				if files := zone.DroppedFiles(); files != nil {
					app.acceptDroppedFiles(files)
				}
				if zone.FileDragOver() {
					zone.Border(2, t.Accent)
				}
				zone.Children(func() {
					ui.Column(c).Gap(4).AlignItems(ui.Center).Children(func() {
						ui.Icon(c, icons.Upload).Size(24, 24).TextColor(t.TextMuted)
						ui.Text(c, "把文件拖到这里")
						ui.Text(c, "支持文件与文件夹").FontSize(12).TextColor(t.TextMuted)
					})
				})
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "已接收文件").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Gap(6).Children(func() {
					if len(d.Files) == 0 {
						ui.Text(c, "还没有收到文件。").TextColor(t.TextMuted)
						return
					}
					for _, f := range d.Files {
						ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
							ui.Text(c, filepath.Base(f.Path)).MaxLines(1).Tooltip(f.Path)
							ui.Text(c, f.meta()).FontSize(12).TextColor(t.TextMuted)
						})
					}
				})

				ui.Text(c, "接收到的路径")
				if ui.Button(c, "清空列表").Disabled(len(d.Files) == 0).
					JustifySelf(ui.Start).Clicked() {
					d.Files = nil
					d.Dropped = 0
				}
			})

		ui.Divider(c)

		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "拖出文件").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, "把下面的示例文件拖到资源管理器或 Finder，会以原生文件数据落到目标应用。")
					if d.SourcePath != "" {
						ui.Text(c, "示例文件: "+d.SourcePath).TextColor(t.TextMuted)
					} else {
						ui.Text(c, d.SourceNote).TextColor(t.TextMuted)
					}
				})

				ui.Text(c, "ui.Element.DragData \rtransfer.FileData")
				if d.SourcePath != "" {
					card := ui.Box(c).Width(280).Height(72).Radius(8).Padding(12, 16).
						Background(t.Surface).Border(1, t.Border).Center().JustifySelf(ui.Start)
					card.DragData(d.Source)
					card.Children(func() {
						ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
							ui.Icon(c, icons.External).Size(24, 24).TextColor(t.Accent)
							ui.Column(c).Gap(2).Children(func() {
								ui.Text(c, filepath.Base(d.SourcePath)).Bold().MaxLines(1)
								ui.Text(c, "按住并拖出").FontSize(12).TextColor(t.TextMuted)
							})
						})
					})
				}
			})
	})
}

// acceptDroppedFiles 记录拖入的路径；文件信息在这里读取一次，视图重绘时不再访问磁盘
func (app *app) acceptDroppedFiles(paths []string) {
	d := &app.pageDragDropData
	for _, path := range paths {
		f := droppedFile{Path: path, Size: -1}
		if info, err := os.Stat(path); err == nil {
			f.Dir = info.IsDir()
			if !f.Dir {
				f.Size = info.Size()
			}
		}
		d.Files = append(d.Files, f)
		d.Dropped++
	}
	if len(d.Files) > maxDroppedFiles {
		d.Files = d.Files[len(d.Files)-maxDroppedFiles:]
	}
}

// loadDragSource 读取应用资源中的图标文件，作为拖出示例的源文件
func (app *app) loadDragSource() {
	d := &app.pageDragDropData
	dir, err := mygo.App.Path(mygo.PathResources)
	if err != nil {
		d.SourceNote = fmt.Sprintf("无法定位资源目录: %v", err)
		return
	}
	path := filepath.Join(dir, "icon.png")
	if _, err := os.Stat(path); err != nil {
		d.SourceNote = fmt.Sprintf("示例文件不可用: %v", err)
		return
	}
	data, err := transfer.FileData(path)
	if err != nil {
		d.SourceNote = fmt.Sprintf("无法创建拖放数据: %v", err)
		return
	}
	d.Source = data
	d.SourcePath = path
}

// meta 返回列表里显示的补充信息
func (f droppedFile) meta() string {
	switch {
	case f.Dir:
		return "文件夹"
	case f.Size < 0:
		return "无法读取大小"
	default:
		return formatTableSize(f.Size)
	}
}
