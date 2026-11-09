package handler

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func call(t *testing.T, method, target, body string) (int, string) {
	t.Helper()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	Handler(w, req)
	b, _ := io.ReadAll(w.Result().Body)
	return w.Code, string(b)
}

func TestRoutes(t *testing.T) {
	for _, tc := range []struct{ method, path, body string }{
		{"GET", "/", ""},
		{"GET", "/about", ""},
		{"GET", "/style.css", ""},
		{"GET", "/app.js", ""},
		{"GET", "/api/risk?rain=250&slope=35&soil=80", ""},
		{"GET", "/api/points", ""},
		{"POST", "/api/risk", `{"rain":300,"slope":40,"soil":90}`},
		{"GET", "/api/index/", ""},
		{"GET", "/api/index/api/points", ""},
	} {
		code, body := call(t, tc.method, tc.path, tc.body)
		if code != 200 {
			t.Fatalf("%s %s → %d (%s)", tc.method, tc.path, code, body[:min(120, len(body))])
		}
		if strings.HasPrefix(tc.path, "/api") && strings.Contains(tc.path, "points") || strings.Contains(tc.path, "risk") {
			var v map[string]any
			if err := json.Unmarshal([]byte(body), &v); err != nil {
				t.Fatalf("%s %s không trả JSON: %v", tc.method, tc.path, err)
			}
		}
	}
	code, body := call(t, "GET", "/khong-ton-tai", "")
	if code != 404 || !strings.Contains(body, "0396241674") {
		t.Fatalf("404 page sai: %d", code)
	}
	_, pb := call(t, "GET", "/api/points", "")
	var pts struct {
		Points []map[string]any `json:"points"`
	}
	if err := json.Unmarshal([]byte(pb), &pts); err != nil || len(pts.Points) != 15 {
		t.Fatalf("cần 15 điểm mock, got %d (err %v)", len(pts.Points), err)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
