package handler_test

import (
	"bytes"
	"encoding/json"
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

	req, _ := http.NewRequest("POST", "/api/urls", bytes.NewBuffer([]byte(request)))
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	var response map[string]map[string]string
	json.Unmarshal([]byte(w.Body.Bytes()), &response)

	assert.Equal(t, "www.google.com", response["url"]["redirectUrl"])
	assert.Equal(t, 4, len(response["url"]["shortenedId"]))
	assert.Equal(t, time.Time{}.UTC().Format(time.RFC3339), response["url"]["expiryDate"])
}
