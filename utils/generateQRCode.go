package utils

import (
	"fmt"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func GenerateQRCode(url string) (string, error) {
	qrc, err := qrcode.New(url)
	if err != nil {
		fmt.Printf("could not generate QRCode:%v", err)
		return "", err
	}
	w, err := standard.New("qrcode.png")
	if err != nil {
		fmt.Printf("standar.New failed:%v", err)
		return "", err
	}

	// save file
	if err = qrc.Save(w); err != nil {
		fmt.Printf("could not save image:%v", err)
		return "", err
	}

	return "qrcode.png", nil
}
