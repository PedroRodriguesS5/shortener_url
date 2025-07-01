package httpclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/pedrorodrigues5/shorter_url/internal/db"
	model "github.com/pedrorodrigues5/shorter_url/models"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/shorten", ShortenerURL)
	r.GET("/:code", ResolveURL)
	return r
}

func TestShortenerURL(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer s.Close()

	os.Setenv("REDIS_URL", s.Addr())
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("DB", "0")
	os.Setenv("BASE_URL", "http://localhost:8080")

	db.InitRedis()

	router := setupRouter()

	t.Run("should shorten a url", func(t *testing.T) {
		url := model.URLMapping{
			URL: "https://google.com",
		}
		jsonValue, _ := json.Marshal(url)
		req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("should resolve a url", func(t *testing.T) {
		code := "test"
		url := "https://google.com"
		db.SaveURL(code, url, 60)

		req, _ := http.NewRequest(http.MethodGet, "/"+code, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusFound {
			t.Errorf("expected status %d, got %d", http.StatusFound, w.Code)
		}
	})
}
