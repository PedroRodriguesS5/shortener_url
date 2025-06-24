package utils

import (
	"github.com/skip2/go-qrcode"
)

// GenerateQRCode generates a QR code for the given URL and returns it as a PNG byte slice.
func GenerateQRCode(url string) ([]byte, error) {
	// Generate the QR code with medium redundancy and 256x256 pixels.
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return png, nil
}
