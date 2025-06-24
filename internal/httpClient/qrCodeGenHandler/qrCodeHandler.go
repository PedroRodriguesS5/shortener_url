package httpclient

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedrorodrigues5/shorter_url/utils"
)

func QRCodeHandler(c *gin.Context) {
	code := c.Param("code")
	pic, err := utils.GenerateQRCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	c.Data(http.StatusOK, "image/png", []byte(pic))
}
