package db

import (
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func TestRedis(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer s.Close()

	os.Setenv("REDIS_URL", s.Addr())
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("DB", "0")

	InitRedis()

	code := "test"
	url := "https://www.google.com"
	expired := 60

	err = SaveURL(code, url, expired)
	if err != nil {
		t.Errorf("SaveURL() error = %v", err)
	}

	resultURL, err := GetURL(code)
	if err != nil {
		t.Errorf("GetURL() error = %v", err)
	}
	if resultURL != url {
		t.Errorf("GetURL() = %v, want %v", resultURL, url)
	}

	err = IncrementClick(code)
	if err != nil {
		t.Errorf("IncrementClick() error = %v", err)
	}

	clicks, err := GetClicks(code)
	if err != nil {
		t.Errorf("GetClicks() error = %v", err)
	}
	if clicks != 1 {
		t.Errorf("GetClicks() = %v, want %v", clicks, 1)
	}
}
