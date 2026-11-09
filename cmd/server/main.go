// Entry cho Vercel Go runtime (server mode) và chạy local.
// Vercel cấp PORT ngẫu nhiên; local mặc định 8080.
// HTK – 0396241674 – cập nhật 09/11/2026
package main

import (
	"log"
	"net/http"
	"os"

	handler "binhdinh-landslide-warning/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", handler.Handler)
	log.Printf("Cảnh báo sạt lở Bình Định (HTK 0396241674) lắng nghe cổng %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}