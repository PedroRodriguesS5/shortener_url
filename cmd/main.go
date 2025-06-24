package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pedrorodrigues5/shorter_url/internal/db"
	httpclientQrCode "github.com/pedrorodrigues5/shorter_url/internal/httpClient/qrCodeGenHandler"
	httpclientShortener "github.com/pedrorodrigues5/shorter_url/internal/httpClient/shortenerUrlHandler"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitRedis()
	r := gin.Default()

	r.POST("/shorten", httpclientShortener.ShortenerURL)
	r.GET("/:code", httpclientShortener.ResolveURL)
	r.GET("/qrcode", httpclientQrCode.QRCodeHandler)
	r.Run(":8080")
}
