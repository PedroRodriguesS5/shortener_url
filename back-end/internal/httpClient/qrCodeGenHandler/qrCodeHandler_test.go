package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestQRCodeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return a qr code", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/qr-code?url=https://google.com", nil)
		w := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/qr-code", QRCodeHandler)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		if w.Header().Get("Content-Type") != "image/png" {
			t.Errorf("expected content type image/png, got %s", w.Header().Get("Content-Type"))
		}
	})

	t.Run("should return an error if no url is provided", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/qr-code", nil)
		w := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/qr-code", QRCodeHandler)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should return an error if an invalid url is provided", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/qr-code?url=invalid-url", nil)
		w := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/qr-code", QRCodeHandler)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
