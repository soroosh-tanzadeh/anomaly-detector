package trafficplot

import (
	"fmt"
	"image/color"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var p *plot.Plot
var traffics plotter.XYs = plotter.XYs{}
var anomalyData plotter.XYs = plotter.XYs{}

func CreateTrafficPlot() {
	p = plot.New()

	p.Title.Text = "InComing/Outcoming Traffic"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Traffic"
}

func Capture(traffic float64) {
	traffics = append(traffics, plotter.XY{X: float64(time.Now().Unix()), Y: traffic})
}

func CaptureAnomaly(anomaly time.Time, traffic float64) {
	anomalyData = append(anomalyData, plotter.XY{X: float64(anomaly.Unix()), Y: traffic})
}

func Save() {
	err := plotutil.AddLinePoints(p,
		"Traffics", traffics,
	)
	if err != nil {
		panic(err)
	}
	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(anomalyData)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 0, B: 255, G: 255, A: 255}
	p.Add(s)
	p.Legend.Add("anomalies", s)
	fmt.Printf("anomaly data %v", anomalyData)
	if err := p.Save(32*vg.Inch, 12*vg.Inch, "result.png"); err != nil {
		panic(err)
	}
}
