package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/asad1123/url-shortener/src/config"
	api "github.com/asad1123/url-shortener/src/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var r *gin.Engine

func TestMain(m *testing.M) {
	config, err := config.LoadConfig("../..", "test")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	r = api.SetupRouter(&config)
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func TestCreateShortenedUrlWithNoExpiry(t *testing.T) {
	w := httptest.NewRecorder()

	request := "{\"redirectUrl\": \"www.google.com\"}"

	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer([]byte(request)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]map[string]string
	json.Unmarshal([]byte(w.Body.Bytes()), &response)

	assert.Equal(t, "www.google.com", response["url"]["redirectUrl"])
	assert.Equal(t, 4, len(response["url"]["shortenedId"]))
	assert.Equal(t, time.Time{}.UTC().Format(time.RFC3339), response["url"]["expiryDate"])
}

func createUrl(longUrl string) map[string]map[string]string {
	w := httptest.NewRecorder()

	request := "{\"redirectUrl\": \"www.google.com\"}"

	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer([]byte(request)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var response map[string]map[string]string
	json.Unmarshal([]byte(w.Body.Bytes()), &response)

	return response
}

func TestCreateShortenedUrlWithExpiry(t *testing.T) {
	w := httptest.NewRecorder()
	request := "{\"redirectUrl\": \"www.google.com\", \"expiryDate\": \"2022-01-01T00:00:00Z\"}"

	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer([]byte(request)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]map[string]string
	json.Unmarshal([]byte(w.Body.Bytes()), &response)

	assert.Equal(t, "www.google.com", response["url"]["redirectUrl"])
	assert.Equal(t, 4, len(response["url"]["shortenedId"]))
	assert.Equal(
		t,
		time.Date(2022, time.January, 01, 0, 0, 0, 0, time.UTC).Format(time.RFC3339),
		response["url"]["expiryDate"],
	)
}

func TestCreateShortenedUrlInvalidRequestError(t *testing.T) {
	w := httptest.NewRecorder()

	request := "{\"blah\": \"www.google.com\"}"

	req, _ := http.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer([]byte(request)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal([]byte(w.Body.Bytes()), &response)

	assert.Equal(t, "Need to provide redirect URL", response["error"])
}

func TestRetrieveUrlToRedirect(t *testing.T) {
	response := createUrl("www.google.com")

	shortenedId := response["url"]["shortenedId"]
	retrieveRoute := fmt.Sprintf("/api/v1/urls/%s", shortenedId)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", retrieveRoute, nil)
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)

	url, _ := w.Result().Location()
	assert.Equal(t, "/api/v1/urls/www.google.com", url.Path)
}

func TestRetrieveUrlToRedirectNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	retrieveRoute := fmt.Sprintf("/api/v1/urls/%s", "abcd")

	req, _ := http.NewRequest("GET", retrieveRoute, nil)
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteUrl(t *testing.T) {
	response := createUrl("www.google.com")

	shortenedId := response["url"]["shortenedId"]
	deleteRoute := fmt.Sprintf("/api/v1/urls/%s", shortenedId)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", deleteRoute, nil)
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteUrlNotFound(t *testing.T) {
	deleteRoute := fmt.Sprintf("/api/v1/urls/%s", "abcd")

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", deleteRoute, nil)
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func retrieveUrl(shortenedId string) {
	retrieveRoute := fmt.Sprintf("/api/v1/urls/%s", shortenedId)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", retrieveRoute, nil)
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)
}

func TestGetUsageAnalyics(t *testing.T) {
	result := createUrl("www.google.com")

	shortenedId := result["url"]["shortenedId"]
	retrieveRoute := fmt.Sprintf("/api/v1/analytics/urls/%s", shortenedId)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", retrieveRoute, nil)
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var response map[string]map[string]interface{}
	json.Unmarshal([]byte(w.Body.Bytes()), &response)

	count := int(response["analytics"]["count"].(float64))
	assert.Equal(t, 0, count)

	retrieveUrl(shortenedId)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", retrieveRoute, nil)
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	json.Unmarshal([]byte(w.Body.Bytes()), &response)

	count = int(response["analytics"]["count"].(float64))
	assert.Equal(t, 1, count)
}
