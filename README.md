# ⛰️ binhdinh-landslide-warning — Cảnh báo sạt lở vùng núi Bình Định

**Cập nhật: 09/11/2026** · Cao điểm **mùa mưa lũ 2026** miền Trung.

## Vấn đề

Tháng 11/2026, sau bão số 9 và số 10, mưa lớn kéo dài khiến các huyện miền núi
**An Lão, Vĩnh Thạnh, Tây Sơn** (Bình Định) đối mặt nguy cơ sạt lở đất nghiêm trọng:
taluy nứt nẻ sát nhà dân, đường liên xã chia cắt, đất đồi trồng keo bão hòa nước.
Dự án cung cấp công cụ ước tính chỉ số nguy cơ sạt lở theo lượng mưa, độ dốc và
độ bão hòa đất, cùng danh sách 15 điểm theo dõi trọng yếu.

> ⚠️ **Dữ liệu mô phỏng – cập nhật 09/11/2026** — không thay thế cảnh báo chính thức của cơ quan chức năng.

## Tính năng

- 🧮 **Go API tính risk index** (0–100): `GET /api/risk?rain=220&slope=35&soil=80` hoặc `POST /api/risk` với JSON `{rain, slope, soil}`.
- 📍 **15 điểm sạt lở mock** tại An Lão, Vĩnh Thạnh, Tây Sơn: `GET /api/points` (JSON, kèm mức + màu + khuyến nghị).
- 📊 Frontend HTML + **Chart.js** gọi trực tiếp API Go, biểu đồ cột tô màu theo 4 mức:
  Thấp (xanh) < 35 · Trung bình (vàng) 35–54 · Cao (cam) 55–74 · Rất cao (đỏ) ≥ 75.
- 🔎 Lọc điểm theo huyện, bảng ghi chú hiện trường, kỹ năng an toàn khi sạt lở.
- 📱 Responsive 375px → 1440px, UI tiếng Việt, trang Giới thiệu + trang 404 có hotline.

## Công nghệ

- **Go 1.24 + Fiber v2 + adaptor v2** (serverless trên Vercel Go runtime, `api/index.go` export `Handler`).
- Frontend tĩnh nhúng bằng `go:embed` — cùng một binary phục vụ cả `/api/*` lẫn giao diện.
- Chart.js 4 (CDN), không cần database.

## Cách chạy local

```bash
go mod download
go build ./...
go test ./...              # smoke test: /, /about, /api/risk, /api/points, 404
go run local/server.go     # mở http://localhost:8080
```

## Deploy

Triển khai Vercel: `vercel --prod --yes` (rewrite toàn bộ path qua `api/index.go`,
middleware khôi phục đường dẫn gốc).

🌐 **Link deploy:** xem [`deploy_url.txt`](deploy_url.txt)

## 📞 Liên hệ HTK – 0396241674 – 09/11/2026

- Tác giả: **HTK**
- Hotline/Zalo: **0396241674** — [zalo.me/0396241674](https://zalo.me/0396241674)
- Email: **htk.contact@gmail.com**
- Phát hiện điểm sạt lở mới hoặc cần hỗ trợ khẩn cấp: gọi/nhắn Zalo ngay.

---

© 2026 HTK · Hotline/Zalo: [0396241674](tel:0396241674) · Email: htk.contact@gmail.com