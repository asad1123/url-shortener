package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/asad1123/url-shortener/src/cache"
	"github.com/asad1123/url-shortener/src/db"
	"github.com/asad1123/url-shortener/src/keygen"
	model_url "github.com/asad1123/url-shortener/src/models/url"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

func CreateShortenedUrl(c *gin.Context) {

	url := model_url.Url{}
	err := c.Bind(&url)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not resolve body"})
	}

	url.CreatedAt = time.Now().UTC()

	// move string length to env variable for configurability
	// this can then be increased if we face higher scale,
	// thus leading to a higher chance of a collision
	url.ShortenedId = keygen.RandomString(4)

	err = saveUrl(url)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new URL."})
	}

	c.JSON(http.StatusOK, gin.H{"url": &url})

}

func saveUrl(url model_url.Url) error {
	err := db.SaveNewUrl(url)
	if err != nil {
		return err
	}

	// write through cache design
	ctx := context.Background()
	err = cache.SaveShortenedUrlToCache(ctx, url.RedirectUrl, url.ShortenedId)

	return err
}

func RetrieveUrlToRedirect(c *gin.Context) {

	id := c.Param("id")
	url, err := getUrl(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find this short URL."})
	}

	urlUsage := model_url.UrlUsage{}
	urlUsage.ShortenedId = id
	urlUsage.AccessedAt = time.Now().UTC()

	err = db.SaveUrlUsage(urlUsage)
	// if analytics fails, we do not want to mark the request as a failure
	// since there is no end user impact
	// however, we should log this as an error on which to trigger actions
	if err != nil {
		msg := fmt.Sprintf("ERROR: Failed to save analytics for shortened URL : %s", id)
		log.Default().Println(msg)
	}

	c.Redirect(http.StatusTemporaryRedirect, url)

}

func getUrl(id string) (string, error) {
	ctx := context.Background()
	redirectUrl, err := cache.GetRedirectUrlFromCache(ctx, id)
	if err != nil {
		// record our cache miss here
		msg := fmt.Sprintf("Cache miss: %s", id)
		log.Println(msg)

		url, err := db.GetUrl(id)

		// save this back to cache for future hits
		cache.SaveShortenedUrlToCache(ctx, url.RedirectUrl, url.ShortenedId)
		return url.RedirectUrl, err
	}

	return redirectUrl, err
}

func DeleteShortenedUrl(c *gin.Context) {

	id := c.Param("id")
	info, err := deleteUrl(id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete this short URL."})
	}
	if info.Removed == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find this short URL."})
	}

	c.Status(http.StatusNoContent)
}

func deleteUrl(id string) (*mgo.ChangeInfo, error) {
	ctx := context.Background()
	err := cache.DeleteRedirectUrlFromCache(ctx, id)
	info, dbErr := db.DeleteUrl(id)

	if err != nil {
		err = errors.Wrap(err, dbErr.Error())
	} else {
		err = dbErr
	}

	return info, err
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
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read analytics for this URL."})
	}

	var response model_url.UrlUsageResponseSchema
	response.ShortenedId = id
	response.Count = count

	c.JSON(http.StatusOK, gin.H{"analytics": &response})
}
