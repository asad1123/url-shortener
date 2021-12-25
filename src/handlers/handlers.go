package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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

	url.CreatedAt = time.Now().UTC()

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

	urlUsage := model_url.UrlUsage{}
	urlUsage.ShortenedId = url.ShortenedId
	urlUsage.AccessedAt = time.Now().UTC()

	err = db.SaveUrlUsage(urlUsage)
	// if analytics fails, we do not want to mark the request as a failure
	// since there is no end user impact
	// however, we should log this as an error on which to trigger actions
	if err != nil {
		msg := fmt.Sprintf("ERROR: Failed to save analytics for shortened URL : %s", url.ShortenedId)
		log.Default().Println(msg)
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

func getInitialTimestamp(td string) (*time.Time, error) {
	// default to zero date
	// this will handle the case where no query parameter
	// was passed as well
	if len(td) == 0 {
		return &time.Time{}, nil
	}

	// ensure time duration is negative
	// it doesn't really make sense to have this positive
	// anyway since all queries would be on historical data
	if !strings.Contains(td, "-") {
		td = "-" + td
	}

	timeDuration, err := time.ParseDuration(td)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().UTC().Add(timeDuration)

	return &timestamp, nil
}

func GetUsageAnalyticsForUrl(c *gin.Context) {

	id := c.Param("id")

	var query model_url.UrlUsageRequestSchema
	c.Bind(&query)

	initialTimestamp, err := getInitialTimestamp(query.Since)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
	}

	count, err := db.SearchUrlUsage(id, *initialTimestamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read analytics for this URL."})
	}

	var response model_url.UrlUsageResponseSchema
	response.ShortenedId = id
	response.Count = count

	c.JSON(http.StatusOK, gin.H{"analytics": &response})
}
