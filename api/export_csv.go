package main

import (
	"encoding/csv"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	data := [][]string{
		{"ID", "Name", "Email"},
		{"1", "Alice", "alice@example.com"},
		{"2", "Bob", "bob@example.com"},
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=users.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	for _, record := range data {
		_ = writer.Write(record)
	}
}
