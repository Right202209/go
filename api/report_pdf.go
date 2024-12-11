package handler

import (
	"net/http"

	"github.com/jung-kurt/gofpdf"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "User Report")

	// Add sample table
	data := [][]string{
		{"ID", "Name", "Email"},
		{"1", "Alice", "alice@example.com"},
		{"2", "Bob", "bob@example.com"},
	}
	for _, line := range data {
		pdf.Ln(10)
		for _, col := range line {
			pdf.Cell(40, 10, col)
		}
	}

	err := pdf.Output(w)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment;filename=report.pdf")
}
