package utils

import (
	"testing"
)

func TestGenerateQRCode(t *testing.T) {
	url := "https://www.google.com"
	png, err := GenerateQRCode(url)
	if err != nil {
		t.Errorf("GenerateQRCode() error = %v", err)
		return
	}
	if len(png) == 0 {
		t.Errorf("GenerateQRCode() png is empty")
	}
}
