package handler

import (
	"net/http"
	"time"

	"github.com/asad1123/url-shortener/src/keygen"
	model_url "github.com/asad1123/url-shortener/src/models/url"
	"github.com/gin-gonic/gin"
)

func CreateShortenedUrl(c *gin.Context) {

	url := model_url.Url{}
	err := c.Bind(&url)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not resolve body"})
	}

	url.CreatedAt = time.Now()

	// move string length to env variable for configurability
	// this can then be increased if we face higher scale,
	// thus leading to a higher chance of a collision
	url.ShortenedId = keygen.RandomString(4)

	c.JSON(http.StatusOK, gin.H{"url": &url})

}

func RetrieveShortenedUrl(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"short": "placeholder"})

}
