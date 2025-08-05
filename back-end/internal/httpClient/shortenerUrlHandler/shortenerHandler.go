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

	expiry := 24 * 60 * 60

	if err := db.SaveURL(code, req.URL, expiry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the url"})
		return
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		// Derive base URL from the incoming request when ENV var is not set.
		scheme := "https"
		if c.Request.TLS == nil {
			scheme = "http"
		}
		// If behind a proxy (e.g., Render / Nginx) honour X-Forwarded-Proto header
		if forwardedProto := c.Request.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
			scheme = forwardedProto
		}

		host := c.Request.Host
		baseURL = scheme + "://" + host
	}
	shortURL := baseURL + "/" + code

	c.JSON(http.StatusOK, model.ShortenResponse{
		ShortURL: shortURL,
	})

}

func ResolveURL(c *gin.Context) {
	code := c.Param("url")

	url, err := db.GetURL(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	db.IncrementClick(code)
	c.Redirect(http.StatusFound, url)
}

// GetURLClicks returns the number of clicks for a shortened URL
func GetURLClicks(c *gin.Context) {
	code := c.Param("code")

	clicks, err := db.GetClicks(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found or no clicks recorded"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   code,
		"clicks": clicks,
	})
}

// GetAllMetrics returns all shortened URLs with their metrics
func GetAllMetrics(c *gin.Context) {
	metrics, err := db.GetAllURLs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve metrics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
		"total":   len(metrics),
	})
}
