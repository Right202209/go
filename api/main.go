package handler

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/export-csv":
		exportCSV(w, r)
	case "/generate-report":
		generateReport(w, r)
	case "/render-chart":
		renderChart(w, r)
	case "/import-file":
		importFile(w, r)
	default:
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
	}
}

// 导出用户数据为 CSV 文件
func exportCSV(w http.ResponseWriter, r *http.Request) {
	data := [][]string{
		{"ID", "Name", "Age"},
		{"1", "Alice", "30"},
		{"2", "Bob", "25"},
		{"3", "Charlie", "35"},
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=users.csv")
	writer := csv.NewWriter(w)
	defer writer.Flush()

	for _, record := range data {
		_ = writer.Write(record)
	}
}

// 动态生成报表内容并返回为纯文本
func generateReport(w http.ResponseWriter, r *http.Request) {
	report := "Dynamic Report\n\n"
	report += "1. Total Users: 3\n"
	report += "2. Average Age: 30\n"
	report += "3. Generated at: " + r.URL.Query().Get("time")

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Disposition", "attachment; filename=report.txt")
	w.Write([]byte(report))
}

// 动态生成图表并返回 PNG 图像
func renderChart(w http.ResponseWriter, r *http.Request) {
	// 创建图像
	width, height := 400, 300
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// 添加简单的条形图
	colors := []color.RGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
	}
	values := []int{100, 200, 150}
	barWidth := 50
	for i, val := range values {
		barHeight := val
		barColor := colors[i%len(colors)]
		x := i*barWidth + 50
		for dx := 0; dx < barWidth; dx++ {
			for dy := 0; dy < barHeight; dy++ {
				img.Set(x+dx, height-50-dy, barColor)
			}
		}
	}

	// 返回图像
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, img)
}

// 上传 CSV 文件并解析数据
func importFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 解析 CSV 文件
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, "Failed to parse CSV: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 返回解析后的数据
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("[\n"))
	for i, record := range records {
		w.Write([]byte(fmt.Sprintf("  {\"ID\": \"%s\", \"Name\": \"%s\", \"Age\": \"%s\"}", record[0], record[1], record[2])))
		if i < len(records)-1 {
			w.Write([]byte(",\n"))
		}
	}
	w.Write([]byte("\n]"))
}
