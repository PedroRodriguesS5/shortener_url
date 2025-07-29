package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pedrorodrigues5/shorter_url/internal/db"
	httpclient "github.com/pedrorodrigues5/shorter_url/internal/httpClient/qrCodeGenHandler"
	httpclientShortener "github.com/pedrorodrigues5/shorter_url/internal/httpClient/shortenerUrlHandler"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	db.InitRedis()
	r := gin.Default()

	// 1. Configure middleware
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowOrigins:     []string{"http://localhost:8080", "https://quick-linker.onrender.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	trustedProxies := strings.Split(os.Getenv("TRUSTED_PROXIES"), ",")
	r.SetTrustedProxies(trustedProxies)

	// 2. Register routes
	r.POST("/shorten", httpclientShortener.ShortenerURL)
	r.GET("/qrcode", httpclient.QRCodeHandler)
	r.GET("/:url", httpclientShortener.ResolveURL)
	r.GET("/stats/:code", httpclientShortener.GetURLClicks)

	// 3. Run the server
	r.Run(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))
}
