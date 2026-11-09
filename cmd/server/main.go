// Chạy local: go run ./cmd/server  → http://localhost:8080
// HTK – 0396241674 – cập nhật 09/11/2026
package main

import (
	"log"
	"net/http"

	handler "binhdinh-landslide-warning/api"
)

func main() {
	http.HandleFunc("/", handler.Handler)
	log.Println("Cảnh báo sạt lở Bình Định (HTK 0396241674) đang chạy tại http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
