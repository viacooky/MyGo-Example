package app

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"
)

type pageTableData struct {
	List     ui.ListState
	Rows     []pageTableRow
	Sort     ui.SortOrder
	lastSort ui.SortOrder
	Selected int
}

type pageTableRow struct {
	Name string
	Size int64
	Kind string
}

func (app *app) pageTable(c *ui.Context) {
	d := &app.pageTableData
	if len(d.Rows) == 0 {
		d.Rows = newTableRows()
		d.Selected = -1
		d.Sort = ui.SortOrder{Column: "Name"}
	}
	d.List.Selected = &d.Selected
	d.List.Sort = &d.Sort
	d.List.Key = func(i int) any { return d.Rows[i].Name }

	app.sortTableRows()

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Table")
	})

	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Grid(c).Margin(PageTitleBarMargin).
			GapX(16).GapY(8).
			ColumnTracks(ui.FitContent(), ui.Fr(1)).
			AlignItems(ui.Baseline).
			Children(func() {
				ui.Text(c, "Table").FontSize(16).Bold()
				ui.Column(c).Margin(PageTitleBarMargin).Children(func() {
					ui.Text(c, fmt.Sprintf("已选择: %s", orNone(app.selectedTableRow())))
					ui.Text(c, fmt.Sprintf("排序: %s %s", orNone(d.Sort.Column), tableDirection(d.Sort.Descending)))
				})

			})
		columns := []ui.TableColumn{
			{Title: "Name", Sortable: true},
			{Title: "Size", Width: 100, Align: ui.End, Sortable: true},
			{Title: "Kind", Width: 100, Sortable: true},
		}
		ui.Table(c, &d.List, columns, len(d.Rows), func(row, col int) {
			r := d.Rows[row]
			switch col {
			case 0:
				ui.Text(c, r.Name).SingleLine()
			case 1:
				ui.Text(c, formatTableSize(r.Size))
			case 2:
				ui.Text(c, r.Kind)
			}
		}).Height(260)
		ui.Text(c, "点击表头可按列排序，拖动表头右边缘可调整列宽。").
			TextColor(c.Theme().TextMuted).Margin(PageTitleBarMargin)
	})
}

// newTableRows 构造表格示例数据；数量足够多以展示滚动效果
func newTableRows() []pageTableRow {
	rows := []pageTableRow{
		{Name: "go.mod", Size: 2 * 1024, Kind: "Go"},
		{Name: "go.sum", Size: 8 * 1024, Kind: "Go"},
		{Name: "main.go", Size: 3 * 1024, Kind: "Go"},
		{Name: "ui_layout.go", Size: 5 * 1024, Kind: "Go"},
		{Name: "README.md", Size: 6 * 1024, Kind: "Markdown"},
		{Name: "LICENSE", Size: 1 * 1024, Kind: "Text"},
		{Name: "icon.png", Size: 350 * 1024, Kind: "Image"},
	}

	kinds := []string{"Go", "JSON", "Image", "PDF", "Text", "Video"}
	exts := []string{"go", "json", "png", "pdf", "txt", "mp4"}
	names := []string{"report", "dataset", "backup", "config", "snapshot", "summary"}
	for i := range 40 {
		k := i % len(kinds)
		rows = append(rows, pageTableRow{
			Name: fmt.Sprintf("attachments/%s_%03d.%s", names[i%len(names)], i+1, exts[k]),
			Size: int64(64+(i*53)%4096) * 1024,
			Kind: kinds[k],
		})
	}
	return rows
}

// sortTableRows 按当前排序应用行顺序；表头点击后排序改变，同时清空选择
func (app *app) sortTableRows() {
	d := &app.pageTableData
	if d.Sort.Column == "" || d.Sort == d.lastSort {
		return
	}
	d.lastSort = d.Sort

	slices.SortStableFunc(d.Rows, func(a, b pageTableRow) int {
		var n int
		switch d.Sort.Column {
		case "Size":
			n = cmp.Compare(a.Size, b.Size)
		case "Kind":
			n = cmp.Compare(a.Kind, b.Kind)
		default:
			n = cmp.Compare(a.Name, b.Name)
		}
		if d.Sort.Descending {
			n = -n
		}
		return n
	})
	d.Selected = -1
}

// selectedTableRow 返回当前选择的行名
func (app *app) selectedTableRow() string {
	d := &app.pageTableData
	if d.Selected < 0 || d.Selected >= len(d.Rows) {
		return ""
	}
	return d.Rows[d.Selected].Name
}

func tableDirection(descending bool) string {
	if descending {
		return "降序"
	}
	return "升序"
}

func formatTableSize(n int64) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.0f KB", float64(n)/(1<<10))
	default:
		return fmt.Sprintf("%d B", n)
	}
}
