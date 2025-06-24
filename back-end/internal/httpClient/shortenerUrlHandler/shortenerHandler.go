package httpclient

import (
	"net/http"
	"os"

	model "github.com/pedrorodrigues5/shorter_url/models"

	"github.com/gin-gonic/gin"
	"github.com/pedrorodrigues5/shorter_url/internal/db"
	"github.com/pedrorodrigues5/shorter_url/utils"
)

func ShortenerURL(c *gin.Context) {
	var req model.URLMapping
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Custom short code already exists"})
		return
	}

	code := req.Custom

	if code == "" {
		code = utils.GenerateShortCode(6)
	}

	// Check if custom code already exists
	if _, err := db.GetURL(code); err == nil {
		c.JSON(400, gin.H{"error": "Custom short code already exists"})
		return
	}

	expiry := req.ExpiresIn
	if expiry == 0 {
		expiry = 3600 * 24 * 7
	}

	if err := db.SaveURL(code, req.URL, expiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the url"})
		return
	}

	c.JSON(http.StatusOK, model.ShortenResponse{
		ShortURL: os.Getenv("BASE_URL") + "/" + code,
	})

}

func ResolveURL(c *gin.Context) {
	code := c.Param("code")

	url, err := db.GetURL(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	db.IncrementClick(code)
	c.Redirect(http.StatusFound, url)
}
