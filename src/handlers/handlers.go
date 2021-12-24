package handler

import (
	"net/http"
	"time"

	"github.com/asad1123/url-shortener/src/db"
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

	err = db.SaveNewUrl(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new URL."})
	}

	c.JSON(http.StatusOK, gin.H{"url": &url})

}

func RetrieveShortenedUrl(c *gin.Context) {

	id := c.Param("id")
	url, err := db.GetUrl(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find this short URL."})
	}

	c.Redirect(http.StatusTemporaryRedirect, url.RedirectUrl)

}

func DeleteShortenedUrl(c *gin.Context) {

	id := c.Param("id")
	info, err := db.DeleteUrl(id)
	if info.Removed == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find this short URL."})
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete this short URL."})
	}

	c.Status(http.StatusNoContent)
}
