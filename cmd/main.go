package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pedrorodrigues5/shorter_url/internal/db"
	"github.com/pedrorodrigues5/shorter_url/internal/httpclient"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitRedis()
	r := gin.Default()
	r.POST("/shorten", httpclient.ShortenerURL)
	r.GET("/:code", httpclient.ResolveURL)
	r.Run(":8080")
}
