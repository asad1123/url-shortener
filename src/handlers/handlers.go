package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/asad1123/url-shortener/src/keygen"
	model_url "github.com/asad1123/url-shortener/src/models/url"
	"github.com/asad1123/url-shortener/src/repository/store"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	store  *store.Store
	keygen *keygen.KeyGen
}

func NewHandler(store *store.Store, keygen *keygen.KeyGen) *Handler {
	return &Handler{
		store,
		keygen,
	}
}

func (h *Handler) Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (h *Handler) CreateShortenedUrl(c *gin.Context) {

	url := model_url.Url{}
	err := c.Bind(&url)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Could not resolve body"})
		return
	}

	if len(url.RedirectUrl) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Need to provide redirect URL"})
		return
	}

	url.CreatedAt = time.Now().UTC()
	url.ShortenedId = h.keygen.RandomString()

	existing, _ := h.store.GetUrl(url.ShortenedId)
	retry := 0
	for len(existing) != 0 && retry < 5 {
		url.ShortenedId = h.keygen.RandomString()
		existing, _ = h.store.GetUrl(url.ShortenedId)

		retry++
		time.Sleep(time.Microsecond) // don't want to overload our server
	}

	if retry == 5 {
		// trigger action from this alert
		log.Println("Error: Failed to generate a unique shortened ID")

		// notify user to try again at a later time
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not generate a short URL. Please try again in a bit."})
		return
	}

	err = h.store.SaveUrl(url)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new URL."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"url": &url})
}

func (h *Handler) RetrieveUrlToRedirect(c *gin.Context) {

	id := c.Param("id")
	url, err := h.store.GetUrl(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Could not find this short URL."})
		return
	}

	urlUsage := model_url.UrlUsage{}
	urlUsage.ShortenedId = id
	urlUsage.AccessedAt = time.Now().UTC()

	err = h.store.SaveUrlUsage(urlUsage)
	// if analytics fails, we do not want to mark the request as a failure
	// since there is no end user impact
	// however, we should log this as an error on which to trigger actions
	if err != nil {
		msg := fmt.Sprintf("ERROR: Failed to save analytics for shortened URL : %s", id)
		log.Default().Println(msg)
	}

	c.Redirect(http.StatusFound, url)
	c.Abort()
}

func (h *Handler) DeleteShortenedUrl(c *gin.Context) {

	id := c.Param("id")
	info, err := h.store.DeleteUrl(id)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not delete this short URL."})
		return
	}
	if info.Removed == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Could not find this short URL."})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetUsageAnalyticsForUrl(c *gin.Context) {

	id := c.Param("id")

	var query model_url.UrlUsageRequestSchema
	c.Bind(&query)

	initialTimestamp, err := getInitialTimestamp(query.Since)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
		return
	}

	count, err := h.store.GetUrlUsage(id, *initialTimestamp)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read analytics for this URL."})
		return
	}

	response := model_url.UrlUsageResponseSchema{ShortenedId: id, Count: count}
	c.JSON(http.StatusOK, gin.H{"analytics": &response})
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
