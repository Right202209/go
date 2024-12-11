package main

import (
	"net/http"

	"github.com/wcharczuk/go-chart/v2"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	graph := chart.BarChart{
		Width:  512,
		Height: 512,
		Bars: []chart.Value{
			{Value: 5, Label: "Category A"},
			{Value: 10, Label: "Category B"},
			{Value: 15, Label: "Category C"},
		},
	}

	w.Header().Set("Content-Type", "image/png")
	err := graph.Render(chart.PNG, w)
	if err != nil {
		http.Error(w, "Failed to render chart", http.StatusInternalServerError)
	}
}
