// Package utils provide some common utilities
package utils

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generatePieItems() []opts.PieData {
	items := make([]opts.PieData, 0)
	for k, v := range PercentileMap {
		items = append(items, opts.PieData{
			Name:  fmt.Sprintf("Schema ID %s", k),
			Value: v,
		})
	}
	return items
}

func createPieChart() *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Schemas Statistics",
			Subtitle: fmt.Sprintf("Snapshot: %s", time.Now().Format(time.RFC822)),
		}),
	)
	pie.AddSeries("pie", generatePieItems()).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show:      opts.Bool(true),
				Formatter: "{b}: {c}%",
			}),
		)
	return pie
}

// GenChart will generate a Pie Chart based on data
func GenChart() {
	page := components.NewPage()
	page.AddCharts(
		createPieChart(),
	)
	f, err := os.Create("pie.html")
	if err != nil {
		panic(err)
	}
	_ = page.Render(io.MultiWriter(f))
}
