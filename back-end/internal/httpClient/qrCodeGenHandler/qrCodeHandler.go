package httpclient

import (
	"net/http"
	"regexp"

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

	// Validate the URL using the provided regex.
	// Note: The regex was corrected to escape the dot in `www.` to `www\.`
	regexPattern := `^((([A-Za-z]{3,9}:(?://)?)(?:[-;:&=+$,\w]+@)?[A-Za-z0-9.-]+|(?:www\.|[-;:&=+$,\w]+@)[A-Za-z0-9.-]+)((?:/[+~%/\.\w-_]*)?\??(?:[-+=&;%@.\w_]*)#?(?:[\w]*))?)$`
	isValid, _ := regexp.MatchString(regexPattern, targetURL)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format provided"})
		return
	}

	pngBytes, err := utils.GenerateQRCode(targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Return the QR code as a PNG image.
	c.Data(http.StatusOK, "image/png", pngBytes)
}
