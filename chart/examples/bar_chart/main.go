package main

//go:generate go run main.go

import (
	"os"

	"github.com/cdvelop/docpdf/canvas"
	"github.com/cdvelop/docpdf/chart"
)

func main() {
	graph := chart.BarChart{
		Title: "Test Bar Chart",
		Background: chart.Style{
			Padding: canvas.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		Bars: []chart.Value{
			{Value: 5.25, Label: "Blue"},
			{Value: 4.88, Label: "Green"},
			{Value: 4.74, Label: "Gray"},
			{Value: 3.22, Label: "Orange"},
			{Value: 3, Label: "Test"},
			{Value: 2.27, Label: "??"},
			{Value: 1, Label: "!!"},
		},
	}

	f, _ := os.Create("output.svg")
	defer f.Close()
	graph.Render(chart.SVG, f)
	// graph.Render(chart.PNG, f)

}
