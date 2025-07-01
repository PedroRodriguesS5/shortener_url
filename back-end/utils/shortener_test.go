package utils

import (
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	length := 6
	code := GenerateShortCode(length)
	if len(code) != length {
		t.Errorf("GenerateShortCode() length = %v, want %v", len(code), length)
	}
}
