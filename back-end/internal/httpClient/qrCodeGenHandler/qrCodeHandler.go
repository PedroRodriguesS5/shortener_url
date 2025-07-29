package httpclient

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/pedrorodrigues5/shorter_url/utils"
)

func QRCodeHandler(c *gin.Context) {
	// Get the URL from the 'url' query parameter.
	targetURL := c.Query("url")
	if targetURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please submit an url"})
		return
	}

	decodedURL, err := url.QueryUnescape(targetURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL encoding"})
		return
	}

	// Validate the URL.
	_, err = url.ParseRequestURI(decodedURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format provided"})
		return
	}

	pngBytes, err := utils.GenerateQRCode(decodedURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Return the QR code as a PNG image.
	c.Header("Access-Control-Allow-Origin", "*")
	c.Data(http.StatusOK, "image/png", pngBytes)
}
