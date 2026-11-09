// Package handler phục vụ API + frontend cảnh báo sạt lở vùng núi Bình Định.
// Dự án: binhdinh-landslide-warning — An Lão, Vĩnh Thạnh, Tây Sơn (tháng 11/2026).
// Tác giả: HTK — Hotline/Zalo: 0396241674 — Email: htk.contact@gmail.com
// Cập nhật: 09/11/2026. Dữ liệu mô phỏng – cập nhật 09/11/2026.
package handler

import (
	"io/fs"
	"math"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"

	app "binhdinh-landslide-warning"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

const dataNote = "Dữ liệu mô phỏng – cập nhật 09/11/2026"

// DiemSatLo là một điểm nguy cơ sạt lở (mock).
type DiemSatLo struct {
	ID         int     `json:"id"`
	Ten        string  `json:"ten"`
	Huyen      string  `json:"huyen"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	Mua24h     float64 `json:"mua24h"` // lượng mưa 24h (mm)
	DoDoc      float64 `json:"doDoc"`  // độ dốc sườn đồi (độ)
	AmDat      float64 `json:"amDat"`  // độ bão hòa đất (%)
	HoDan      int     `json:"hoDan"`  // số hộ dân bị đe dọa
	GhiChu     string  `json:"ghiChu"`
	ChiSo      float64 `json:"chiSo"`      // risk index 0–100
	MucDo      string  `json:"mucDo"`      // mức cảnh báo
	Mau        string  `json:"mau"`        // màu hiển thị
	KhuyenNghi string  `json:"khuyenNghi"` // khuyến nghị hành động
}

// rawDiem là dữ liệu gốc 15 điểm mock (chỉ số được tính khi serve).
type rawDiem struct {
	ID     int
	Ten    string
	Huyen  string
	Lat    float64
	Lng    float64
	Mua24h float64
	DoDoc  float64
	AmDat  float64
	HoDan  int
	GhiChu string
}

var diemSatLo = []rawDiem{
	{1, "Thôn An Toàn (An Lão)", "An Lão", 14.585, 108.658, 312, 38, 92, 46, "Vết nứt taluy dài 40m sau mưa lớn 08/11/2026"},
	{2, "Xã An Hưng – suối Đắk Long", "An Lão", 14.552, 108.702, 268, 33, 85, 31, "Đất bão hòa, nguy cơ lũ quét lòng suối"},
	{3, "Xã An Trung – đèo An Khê", "An Lão", 14.498, 108.741, 240, 41, 80, 58, "Tuyến liên xã bị chia cắt cục bộ"},
	{4, "Thôn Gò Cây Duối (An Vinh)", "An Lão", 14.531, 108.812, 198, 28, 74, 22, "Sườn đồi trồng keo bị xói mạnh"},
	{5, "Xã An Toàn – điểm trường cũ", "An Lão", 14.601, 108.640, 176, 24, 68, 15, "Theo dõi vết nứt tường rào"},
	{6, "Thôn Định Bình (Vĩnh Thạnh)", "Vĩnh Thạnh", 14.372, 108.729, 295, 36, 90, 52, "Sạt taluy đường Hồ Chí Nhánh Đông"},
	{7, "Xã Vĩnh Sơn – suối Kôn thượng", "Vĩnh Thạnh", 14.428, 108.611, 254, 40, 86, 37, "Nguy cơ sạt lở bờ suối khi lũ lên"},
	{8, "Xã Vĩnh Kim – đồi keo H12", "Vĩnh Thạnh", 14.335, 108.688, 214, 31, 78, 26, "Đất dốc trồng keo, rễ bám kém"},
	{9, "Thôn Vĩnh Hà (Vĩnh Thịnh)", "Vĩnh Thạnh", 14.302, 108.775, 182, 26, 70, 19, "Mương thoát nước bị bồi lấp"},
	{10, "Xã Vĩnh Thạnh – QL19 đoạn km42", "Vĩnh Thạnh", 14.258, 108.802, 160, 35, 64, 12, "Đá rơi lẻ tẻ khi mưa kéo dài"},
	{11, "Đèo Vĩnh Thạnh – Phú Phong", "Tây Sơn", 14.121, 108.842, 232, 42, 82, 0, "Nguy cơ đá lở chắn đường, cấm dừng đỗ"},
	{12, "Xã Bình Nghi – đồi Trà Dơn", "Tây Sơn", 14.089, 108.905, 205, 29, 76, 34, "Khu dân cư chân đồi cần sơ tán khi mưa >200mm"},
	{13, "Thôn Thuận Ninh (Bình Tân)", "Tây Sơn", 14.154, 108.878, 172, 23, 66, 21, "Vết nứt sân sau 3 hộ dân"},
	{14, "Xã Tây Thuận – kênh chính Phú Phong", "Tây Sơn", 14.065, 108.861, 148, 19, 60, 17, "Sạt mái kênh, theo dõi khi xả lũ"},
	{15, "Xã Bình Tường – đồi Gò Găng", "Tây Sơn", 14.112, 108.923, 126, 17, 55, 9, "Nguy cơ thấp, kiểm tra định kỳ"},
}

// ketQua đánh giá chỉ số rủi ro 0–100 theo lượng mưa, độ dốc, độ ẩm đất.
// Công thức mô phỏng phục vụ minh họa, không thay thế cảnh báo chính thức.
func danhGia(mua, doc, am float64) (chiSo float64, muc, mau, kn string) {
	m := math.Min(math.Max(mua, 0), 400) / 400
	d := math.Min(math.Max(doc, 0), 60) / 60
	a := math.Min(math.Max(am, 0), 100) / 100
	chiSo = math.Round((0.5*m+0.3*d+0.2*a)*1000) / 10
	switch {
	case chiSo >= 75:
		return chiSo, "Rất cao", "#dc2626", "Sơ tán khẩn cấp các hộ chân đồi, cấm qua lại khu vực sạt lở."
	case chiSo >= 55:
		return chiSo, "Cao", "#ea580c", "Chuẩn bị sơ tán, trực canh 24/24, theo dõi vết nứt và nước đục."
	case chiSo >= 35:
		return chiSo, "Trung bình", "#d97706", "Theo dõi mưa từng giờ, khơi thông mương thoát nước."
	default:
		return chiSo, "Thấp", "#16a34a", "Sinh hoạt bình thường, kiểm tra taluy sau mỗi đợt mưa."
	}
}

func tinhDiem() []DiemSatLo {
	out := make([]DiemSatLo, len(diemSatLo))
	for i, r := range diemSatLo {
		d := DiemSatLo{ID: r.ID, Ten: r.Ten, Huyen: r.Huyen, Lat: r.Lat, Lng: r.Lng,
			Mua24h: r.Mua24h, DoDoc: r.DoDoc, AmDat: r.AmDat, HoDan: r.HoDan, GhiChu: r.GhiChu}
		d.ChiSo, d.MucDo, d.Mau, d.KhuyenNghi = danhGia(d.Mua24h, d.DoDoc, d.AmDat)
		out[i] = d
	}
	return out
}

var (
	once     sync.Once
	fiberApp *fiber.App
	staticFS fs.FS
)

func buildApp() *fiber.App {
	a := fiber.New(fiber.Config{DisableStartupMessage: true})
	staticFS, _ = fs.Sub(app.PublicFS, "public")

	api := a.Group("/api")
	api.Get("/points", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"ok": true, "note": dataNote, "total": len(diemSatLo), "points": tinhDiem()})
	})
	api.Get("/risk", func(c *fiber.Ctx) error {
		mua, e1 := strconv.ParseFloat(c.Query("rain", "0"), 64)
		doc, e2 := strconv.ParseFloat(c.Query("slope", "0"), 64)
		am, e3 := strconv.ParseFloat(c.Query("soil", "50"), 64)
		if e1 != nil || e2 != nil || e3 != nil {
			return c.Status(400).JSON(fiber.Map{"ok": false, "error": "Tham số rain/slope/soil phải là số."})
		}
		chiSo, muc, mau, kn := danhGia(mua, doc, am)
		return c.JSON(fiber.Map{"ok": true, "note": dataNote, "input": fiber.Map{"rain": mua, "slope": doc, "soil": am},
			"riskIndex": chiSo, "level": muc, "color": mau, "advice": kn})
	})
	api.Post("/risk", func(c *fiber.Ctx) error {
		var body struct {
			Rain  float64 `json:"rain"`
			Slope float64 `json:"slope"`
			Soil  float64 `json:"soil"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"ok": false, "error": "JSON không hợp lệ. Cần {rain, slope, soil}."})
		}
		chiSo, muc, mau, kn := danhGia(body.Rain, body.Slope, body.Soil)
		return c.JSON(fiber.Map{"ok": true, "note": dataNote, "input": body,
			"riskIndex": chiSo, "level": muc, "color": mau, "advice": kn})
	})

	serveFile := func(c *fiber.Ctx, name string, status int) error {
		b, err := fs.ReadFile(staticFS, name)
		if err != nil {
			return c.Status(500).SendString("Không đọc được tài nguyên.")
		}
		switch path.Ext(name) {
		case ".html":
			c.Type("html", "utf-8")
		case ".css":
			c.Type("css", "utf-8")
		case ".js":
			c.Type("js", "utf-8")
		}
		return c.Status(status).Send(b)
	}

	a.Get("/", func(c *fiber.Ctx) error { return serveFile(c, "index.html", 200) })
	a.Get("/about", func(c *fiber.Ctx) error { return serveFile(c, "about.html", 200) })
	a.Get("/404", func(c *fiber.Ctx) error { return serveFile(c, "404.html", 404) })
	a.Get("/*", func(c *fiber.Ctx) error {
		p := strings.TrimPrefix(c.Params("*"), "/")
		if p == "" {
			p = "index.html"
		}
		if strings.HasPrefix(p, "api/") {
			return c.Status(404).JSON(fiber.Map{"ok": false, "error": "Không tìm thấy API."})
		}
		if f, err := fs.ReadFile(staticFS, p); err == nil {
			switch path.Ext(p) {
			case ".html":
				c.Type("html", "utf-8")
			case ".css":
				c.Type("css", "utf-8")
			case ".js":
				c.Type("js", "utf-8")
			}
			return c.Send(f)
		}
		return serveFile(c, "404.html", 404)
	})
	return a
}

// restorePath khôi phục đường dẫn gốc khi Vercel rewrite qua /api/index/...
func restorePath(r *http.Request) {
	p := r.URL.Path
	if orig := r.Header.Get("X-Original-Path"); orig != "" {
		p = orig
	} else {
		for _, prefix := range []string{"/api/index.go", "/api/index"} {
			if p == prefix {
				p = "/"
				break
			}
			if strings.HasPrefix(p, prefix+"/") {
				p = strings.TrimPrefix(p, prefix)
				break
			}
		}
	}
	r.URL.Path = p
	// adaptor/fasthttp đọc RequestURI: phải đồng bộ cả đường dẫn + query.
	if q := r.URL.RawQuery; q != "" {
		r.RequestURI = p + "?" + q
	} else {
		r.RequestURI = p
	}
}

// Handler là entrypoint cho Vercel Go runtime (api/index.go).
func Handler(w http.ResponseWriter, r *http.Request) {
	restorePath(r)
	once.Do(func() { fiberApp = buildApp() })
	adaptor.FiberApp(fiberApp).ServeHTTP(w, r)
}
