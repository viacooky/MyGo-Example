package app

import (
	"fmt"
	"my-app/internal/icons"

	"github.com/egoist/mygo/ui"
)

type pageListData struct {
	List   ui.ListState
	chosen ui.Selection[string]
	Files  []pageListFile
}

type pageListFile struct {
	Name string
	Size int64
}

func (app *app) pageList(c *ui.Context) {
	app.pageListData.List.Key = func(i int) any { return app.pageListData.Files[i].Name }
	app.pageListData.List.Selection = &app.pageListData.chosen

	if len(app.pageListData.Files) == 0 {
		app.pageListData.Files = []pageListFile{
			{Name: "empty_file.dat", Size: 0},                       // 测试 0 字节
			{Name: "config.yaml", Size: 20 * 1024},                  // 0.02 MB
			{Name: "main.go", Size: 50 * 1024},                      // 0.05 MB
			{Name: "go.sum", Size: 120 * 1024},                      // 0.12 MB
			{Name: "image_001.png", Size: 350 * 1024},               // 0.34 MB
			{Name: "package-lock.json", Size: 850 * 1024},           // 0.83 MB
			{Name: "audio_record.mp3", Size: 5 * 1024 * 1024},       // 5.00 MB
			{Name: "presentation.pptx", Size: 25 * 1024 * 1024},     // 25.00 MB
			{Name: "dataset_v1.csv", Size: 125 * 1024 * 1024},       // 125.00 MB
			{Name: "node_modules.zip", Size: 450 * 1024 * 1024},     // 450.00 MB
			{Name: "video_tutorial.mp4", Size: 1024 * 1024 * 1024},  // 1024.00 MB
			{Name: "backup_database.sql", Size: 2048 * 1024 * 1024}, // 2048.00 MB
			{Name: "Linux_Image.iso", Size: 4096 * 1024 * 1024},     // 4096.00 MB
		}
	}

	ui.Row(c).Children(func() {
		ui.Text(c, "列表").FontSize(PageTitleFontSize).Margin(PageTitleBarMargin).Bold()
	})

	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "List").FontSize(16).Bold()
				ui.List(c, &app.pageListData.List, len(app.pageListData.Files)).Rows(func(row ui.ListRow) {
					isSelected := app.pageListData.List.CellSelected(row.Index, 0)
					var bgColor ui.Color
					var textColor ui.Color
					if isSelected {
						// 选中状态：使用主题强调色作为背景，反转文字颜色
						bgColor = c.Theme().Accent
						textColor = c.Theme().AccentText
					} else {
						// 未选中状态：保持透明背景，使用普通文字颜色
						bgColor = ui.Transparent
						textColor = c.Theme().Text
					}
					ui.Row(row.Context).
						Background(bgColor).
						Padding(8, 16). // 上下 8，左右 16
						FillWidth().
						GapX(8).
						Children(func() {
							size := fmt.Sprintf("%.2f MB", float64(app.pageListData.Files[row.Index].Size)/1024/1024)
							ui.Icon(c, icons.File).Size(24, 24)
							ui.Text(row.Context, app.pageListData.Files[row.Index].Name).TextColor(textColor)
							ui.Spacer(row.Context)
							ui.Text(row.Context, size).TextColor(textColor).TextAlign(ui.End)
						})
				}).Children(func() {
					if len(app.pageListData.Files) == 0 {
						ui.Text(c, "No Files").TextColor(c.Theme().TextMuted)
					}
				})
			})

		ui.Divider(c)
		ui.Column(c).Padding(16).Gap(8).Children(func() {
			ui.Text(c, "支持多选，按 Ctrl 或 Shift 键选中多个文件，Ctrl + A 全选").TextColor(ui.Hex("#FF0000"))
			ui.Text(c, fmt.Sprintf("已选择 %d 项：", app.pageListData.chosen.Len())).Bold()
			for _, file := range app.pageListData.Files {
				if app.pageListData.chosen.Has(file.Name) {
					ui.Text(c, "  • "+file.Name)
				}
			}
		})
	})
}
