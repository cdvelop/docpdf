package main

//go:generate go run main.go

import (
	"os"

	"github.com/cdvelop/docpdf/chart"
	"github.com/cdvelop/docpdf/chart/roboto"
)

func main() {
	// Creamos un nuevo motor de gráficos con la fuente Roboto
	engine, err := chart.NewEngine(roboto.Roboto)
	if err != nil {
		panic(err)
	}

	// Creamos un gráfico de dona utilizando el motor
	values := []chart.Value{
		{Value: 5, Label: "Blue"},
		{Value: 5, Label: "Green"},
		{Value: 4, Label: "Gray"},
		{Value: 4, Label: "Orange"},
		{Value: 3, Label: "Deep Blue"},
		{Value: 3, Label: "test"},
	}

	// Creamos el gráfico con el motor (sintaxis encadenada)
	donutChart := engine.DonutChart(values)
	// Renderizamos a un archivo PNG
	f, _ := os.Create("output.png")
	defer f.Close()
	donutChart.Render(chart.PNG, f)
}
