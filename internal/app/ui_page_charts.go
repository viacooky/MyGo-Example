package app

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/chart"
)

// 图表示例使用的配色与固定维度
var (
	chartPalette = []ui.Color{
		ui.Hex("#3b82f6"),
		ui.Hex("#10b981"),
		ui.Hex("#f59e0b"),
		ui.Hex("#ec4899"),
	}

	chartMonths = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	chartSources   = []string{"Search", "Social", "Direct"}
	chartPlatforms = []string{"macOS", "Windows", "Linux", "Browsers"}
)

// 各图表的样本数量
const (
	chartVisitsLen  = 12
	chartLineLen    = 40
	chartScatterLen = 10
)

// chartPoint 是散点图上的一个点
type chartPoint struct{ X, Y float64 }

type pageChartsData struct {
	Sales        []float64
	Month        int // 选中的柱子，-1 表示未选中
	Line         []float64
	LineSelected int
	Visits       [][]float64
	Stacking     int // 0 为数值，1 为占比
	Scatter      [3][]chartPoint
	Shares       []float64
	initialized  bool
}

func (app *app) pageCharts(c *ui.Context) {
	d := &app.pageChartsData
	if !d.initialized {
		d.initialized = true
		d.Month = -1
		d.LineSelected = -1
		d.Sales = randomChartSales()
		d.Line = randomChartLine(chartLineLen)
		d.Visits = randomChartVisits(len(chartSources), chartVisitsLen)
		for k := range d.Scatter {
			d.Scatter[k] = randomChartScatter(k)
		}
		d.Shares = randomChartShares()
	}

	ui.Row(c).Children(func() {
		app.breadcrumbTitle(c, "Chart")
	})

	ui.Scroll(c).Fill().Padding(0, 16, 64, 16).ScrollbarInsets(0, 0, 20).Children(func() {
		ui.Text(c, "ui/chart 提供坐标轴、比例尺与图形标记，外观由使用方定义。").
			TextColor(c.Theme().TextMuted).Margin(PageTitleBarMargin)
		ui.Row(c).Gap(16).Wrap().AlignItems(ui.Start).Children(func() {
			chartCard(c, "柱状图 · chart.Rect", app.refillChartSales, func() { app.chartBars(c) })
			chartCard(c, "折线与面积 · chart.Area", app.refillChartLine, func() { app.chartLine(c) })
			chartCard(c, "堆叠面积 · chart.Stack", app.refillChartVisits, func() { app.chartStacked(c) })
			chartCard(c, "散点图 · chart.Point", app.refillChartScatter, func() { app.chartScatter(c) })
			chartCard(c, "环形图 · chart.Pie", app.refillChartShares, func() { app.chartDonut(c) })
		})
	})
}

// chartCard 是图表示例的卡片外壳，右上角的「换一批」重新生成该图表的数据
func chartCard(c *ui.Context, title string, refresh func(), build func()) {
	t := c.Theme()
	ui.Column(c).Basis(330).Grow(1).MinWidth(0).Padding(16).Radius(10).
		Background(t.Surface).Border(1, t.Border).Gap(10).Children(func() {
		ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
			ui.Text(c, title).Bold().Grow(1)
			if ui.Button(c, "换一批").Clicked() {
				refresh()
			}
		})
		build()
	})
}

// chartAxisLabel 是坐标轴标签的样式
func chartAxisLabel(c *ui.Context) func(chart.Side, chart.Tick) {
	return func(_ chart.Side, tick chart.Tick) {
		ui.Text(c, tick.Label).NoWrap().FontSize(11).TextColor(c.Theme().TextMuted)
	}
}

// chartGridLines 在图表区域按刻度画出网格线
func chartGridLines(area ui.Element, x, y []chart.Tick, color ui.Color) {
	area.Draw(func(p *ui.Painter, r ui.Rect) {
		for _, tick := range y {
			p.Line(r.X, chart.Y(r, tick.At), r.X+r.W, chart.Y(r, tick.At), 1, color)
		}
		for _, tick := range x {
			p.Line(chart.X(r, tick.At), r.Y, chart.X(r, tick.At), r.Y+r.H, 1, color)
		}
	})
}

// chartLegend 展示各系列的名称与颜色
func chartLegend(c *ui.Context, names []string) {
	ui.Row(c).Gap(14).Wrap().Children(func() {
		for i, name := range names {
			ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(10, 10).Radius(3).Background(chartPalette[i])
				ui.Text(c, name).FontSize(12).TextColor(c.Theme().TextMuted)
			})
		}
	})
}

// chartBars 是可点击选择的柱状图
func (app *app) chartBars(c *ui.Context) {
	d := &app.pageChartsData
	t := c.Theme()

	text := "点击柱子选择月份"
	if d.Month >= 0 {
		text = fmt.Sprintf("%s: $%.1fk", chartMonths[d.Month], d.Sales[d.Month])
	}
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		ui.Text(c, text).TextColor(t.TextMuted).Grow(1)
	})

	x := chart.NewBand(chartMonths...).Padding(0.3, 0.1)
	_, hi := chart.Extent(d.Sales...)
	y := chart.NewLinear(0, hi).Nice(4)
	y.Format = func(v float64) string { return fmt.Sprintf("$%.0fk", v) }
	ticks := y.Ticks(4)

	// 隔一个显示一个月，避免标签重叠
	var months []chart.Tick
	for i, tick := range x.Ticks() {
		if i%2 == 0 {
			months = append(months, tick)
		}
	}

	chart.Plot(c, chart.Axes{Left: ticks, Bottom: months, Label: chartAxisLabel(c)}, func(area ui.Element) {
		chartGridLines(area, nil, ticks, t.Border)
		for i, v := range d.Sales {
			bar := chart.Rect(c.Key(i), x.Band(i), y.Span(0, v)).Radius(3).Background(t.Accent).
				Tooltip(fmt.Sprintf("%s: $%.1fk", chartMonths[i], v)).Label(chartMonths[i])
			if d.Month >= 0 && i != d.Month {
				bar.Background(t.Accent.Alpha(0.35))
			}
			if bar.Clicked() {
				if d.Month == i {
					d.Month = -1
				} else {
					d.Month = i
				}
			}
		}
		if i := d.Month; i >= 0 {
			chart.Anchor(c, x.At(i), y.At(d.Sales[i]), ui.Center, ui.End, func() {
				ui.Text(c, fmt.Sprintf("%.0f", d.Sales[i])).NoWrap().FontSize(11).Bold().Margin(0, 0, 4, 0)
			})
		}
	}).Label("按月销售").Height(190).GapX(8).GapY(6)
}

// chartLine 是带悬停选择的折线面积图
func (app *app) chartLine(c *ui.Context) {
	d := &app.pageChartsData
	t := c.Theme()
	n := len(d.Line)

	x := chart.NewLinear(0, float64(n-1))
	y := chart.NewLinear(0, 1)
	y.Format = func(v float64) string { return fmt.Sprintf("%.0f%%", 100*v) }
	xt, yt := x.Ticks(4), y.Ticks(4)
	xs := make([]float32, n)
	ys := make([]float32, n)
	for i, v := range d.Line {
		xs[i], ys[i] = x.At(float64(i)), y.At(v)
	}
	describe := func(i int) string { return fmt.Sprintf("#%d: %.0f%%", i+1, 100*d.Line[i]) }

	text := "把指针移到折线上选择数据点"
	if d.LineSelected >= 0 && d.LineSelected < n {
		text = describe(d.LineSelected)
	}
	ui.Text(c, text).TextColor(t.TextMuted)

	chart.Plot(c, chart.Axes{Left: yt, Bottom: xt, Label: chartAxisLabel(c)}, func(area ui.Element) {
		area.Draw(func(p *ui.Painter, r ui.Rect) {
			for _, tick := range yt {
				p.Line(r.X, chart.Y(r, tick.At), r.X+r.W, chart.Y(r, tick.At), 1, t.Border)
			}
			p.FillPath(chart.Area(r, xs, nil, ys, chart.CurveMonotone), t.Accent.Alpha(0.12))
			p.StrokePath(chart.Line(r, xs, ys, chart.CurveMonotone), 2, t.Accent)
		})
		chart.Select(c, area, &d.LineSelected, xs, describe)
		if i := d.LineSelected; i >= 0 && i < n {
			chart.RuleX(c, xs[i], 1).Background(t.TextMuted)
			chart.Point(c, xs[i], ys[i], 10, 10).Radius(5).Background(t.Accent).Border(2, t.Background)
		}
	}).Label("样本折线").Height(190).GapX(8).GapY(6)
}

// chartStacked 是按数值或占比堆叠的面积图
func (app *app) chartStacked(c *ui.Context) {
	d := &app.pageChartsData
	t := c.Theme()

	ui.Segmented(c, &d.Stacking, "数值", "占比").Label("堆叠方式").AlignSelf(ui.Start)
	how := chart.StackStandard
	if d.Stacking == 1 {
		how = chart.StackNormalized
	}
	stacked := chart.Stack(d.Visits, how)
	hi := 0.0
	for _, seg := range stacked[len(stacked)-1] {
		hi = max(hi, seg.Hi)
	}
	y := chart.NewLinear(0, hi).Nice(4)
	if how == chart.StackNormalized {
		y = chart.NewLinear(0, 1)
		y.Format = func(v float64) string { return fmt.Sprintf("%.0f%%", 100*v) }
	}
	weeks := len(d.Visits[0])
	x := chart.NewLinear(1, float64(weeks))
	x.Format = func(v float64) string { return fmt.Sprintf("W%.0f", v) }
	xt, yt := x.Ticks(6), y.Ticks(4)
	xs := make([]float32, weeks)
	for i := range xs {
		xs[i] = x.At(float64(i + 1))
	}

	chart.Plot(c, chart.Axes{Left: yt, Bottom: xt, Label: chartAxisLabel(c)}, func(area ui.Element) {
		area.Draw(func(p *ui.Painter, r ui.Rect) {
			for _, tick := range yt {
				p.Line(r.X, chart.Y(r, tick.At), r.X+r.W, chart.Y(r, tick.At), 1, t.Border)
			}
			for s, series := range stacked {
				lo, hi := make([]float32, weeks), make([]float32, weeks)
				for i, seg := range series {
					lo[i], hi[i] = y.At(seg.Lo), y.At(seg.Hi)
				}
				p.FillPath(chart.Area(r, xs, lo, hi, chart.CurveMonotone), chartPalette[s].Alpha(0.75))
				p.StrokePath(chart.Line(r, xs, hi, chart.CurveMonotone), 1.5, chartPalette[s])
			}
		})
	}).Label("按来源的访问量").Height(170).GapX(8).GapY(6)
	chartLegend(c, chartSources)
}

// chartScatter 是按类别分组的散点图
func (app *app) chartScatter(c *ui.Context) {
	d := &app.pageChartsData
	t := c.Theme()

	chartLegend(c, []string{"小型", "中型", "大型"})
	var weights, uses []float64
	for _, class := range d.Scatter {
		for _, p := range class {
			weights = append(weights, p.X)
			uses = append(uses, p.Y)
		}
	}
	x := chart.NewLinear(chart.Extent(weights...)).Nice(4)
	y := chart.NewLinear(chart.Extent(uses...)).Nice(4)
	x.Format = func(v float64) string { return fmt.Sprintf("%.1f t", v/1000) }
	y.Format = func(v float64) string { return fmt.Sprintf("%.0f L", v) }
	xt, yt := x.Ticks(4), y.Ticks(4)

	chart.Plot(c, chart.Axes{Left: yt, Bottom: xt, Label: chartAxisLabel(c)}, func(area ui.Element) {
		chartGridLines(area, xt, yt, t.Border)
		for k, class := range d.Scatter {
			for _, p := range class {
				chart.Point(c, x.At(p.X), y.At(p.Y), 9, 9).Radius(4.5).
					Background(chartPalette[k].Alpha(0.8)).Border(1, chartPalette[k]).
					Tooltip(fmt.Sprintf("%.0f kg, %.1f L/100 km", p.X, p.Y))
			}
		}
	}).Label("车重与油耗").Height(190).GapX(8).GapY(6)
}

// chartDonut 是悬停突出扇区的环形图
func (app *app) chartDonut(c *ui.Context) {
	d := &app.pageChartsData
	t := c.Theme()
	slices := chart.Pie(d.Shares, 0.025)
	const size, inner, outer = 170, 52, 80

	ui.Row(c).Gap(20).AlignItems(ui.Center).Children(func() {
		ring := ui.Box(c).Size(size, size).Label("平台占比")
		hover := -1
		if x, y, over := ring.PointerPosition(); over {
			if i, ok := chart.SliceAt(slices, x-size/2, y-size/2, inner, outer); ok {
				hover = i
			}
		}
		ring.Draw(func(p *ui.Painter, r ui.Rect) {
			for i, s := range slices {
				out := float32(outer - 4)
				if i == hover {
					out = outer
				}
				p.FillPath(chart.Sector(r.X+size/2, r.Y+size/2, inner, out, s), chartPalette[i])
			}
		})
		ring.Children(func() {
			ui.Column(c).Fill().Center().Children(func() {
				if hover < 0 {
					ui.Text(c, "4.2M").FontSize(22).Bold()
					ui.Text(c, "用户").FontSize(12).TextColor(t.TextMuted)
					return
				}
				ui.Text(c, fmt.Sprintf("%.0f%%", d.Shares[hover])).FontSize(22).Bold()
				ui.Text(c, chartPlatforms[hover]).FontSize(12).TextColor(t.TextMuted)
			})
		})

		ui.Column(c).Gap(8).Children(func() {
			for i, name := range chartPlatforms {
				ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(10, 10).Radius(3).Background(chartPalette[i])
					ui.Text(c, name).FontSize(13).Width(70)
					ui.Text(c, fmt.Sprintf("%.0f%%", d.Shares[i])).FontSize(13).TextColor(t.TextMuted)
				})
			}
		})
	})
}

// refillChartSales 重新生成柱状图数据
func (app *app) refillChartSales() {
	d := &app.pageChartsData
	d.Sales = randomChartSales()
	d.Month = -1
}

// refillChartLine 重新生成折线数据
func (app *app) refillChartLine() {
	d := &app.pageChartsData
	d.Line = randomChartLine(len(d.Line))
	d.LineSelected = -1
}

// refillChartVisits 重新生成堆叠面积的数据
func (app *app) refillChartVisits() {
	d := &app.pageChartsData
	d.Visits = randomChartVisits(len(d.Visits), len(d.Visits[0]))
}

// refillChartScatter 重新生成各组散点
func (app *app) refillChartScatter() {
	d := &app.pageChartsData
	for k := range d.Scatter {
		d.Scatter[k] = randomChartScatter(k)
	}
}

// refillChartShares 重新生成环形图的占比
func (app *app) refillChartShares() {
	d := &app.pageChartsData
	d.Shares = randomChartShares()
}

// randomChartSales 生成带趋势波动的月度销售额
func randomChartSales() []float64 {
	phase := 2 * math.Pi * rand.Float64()
	values := make([]float64, len(chartMonths))
	for i := range values {
		values[i] = 40 + 40*math.Sin(float64(i)/3.5+phase) + 20*rand.Float64()
	}
	return values
}

// randomChartLine 生成光滑起伏的折线样本
func randomChartLine(n int) []float64 {
	phase := 2 * math.Pi * rand.Float64()
	freq := 2 + 3*rand.Float64()
	values := make([]float64, n)
	for i := range values {
		values[i] = 0.5 + 0.3*math.Sin(phase+float64(i)*freq/float64(n)) + 0.15*rand.Float64()
	}
	return values
}

// randomChartVisits 生成按来源分组的每周访问量
func randomChartVisits(series, weeks int) [][]float64 {
	values := make([][]float64, series)
	for s := range values {
		values[s] = make([]float64, weeks)
		base := 15 + 40*rand.Float64()
		for i := range values[s] {
			values[s][i] = base + 40*rand.Float64() + 20*float64(i)/float64(weeks)
		}
	}
	return values
}

// randomChartScatter 生成第 class 组散点
func randomChartScatter(class int) []chartPoint {
	points := make([]chartPoint, chartScatterLen)
	for i := range points {
		x := 1100 + 450*float64(class) + 170*rand.NormFloat64()
		points[i] = chartPoint{X: x, Y: max(1, 3+x/420+0.9*rand.NormFloat64())}
	}
	return points
}

// randomChartShares 生成总和为 100 的整数占比
func randomChartShares() []float64 {
	weights := make([]float64, len(chartPlatforms))
	sum := 0.0
	for i := range weights {
		weights[i] = 5 + 60*rand.Float64()
		sum += weights[i]
	}
	shares := make([]float64, len(weights))
	total := 0.0
	for i, w := range weights {
		shares[i] = math.Round(w / sum * 100)
		total += shares[i]
	}
	shares[len(shares)-1] += 100 - total // 舍入误差补到最后一项
	return shares
}
