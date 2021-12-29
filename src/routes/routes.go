package routes

import (
	"fmt"
	"log"

	"github.com/asad1123/url-shortener/src/config"
	handlers "github.com/asad1123/url-shortener/src/handlers"
	"github.com/asad1123/url-shortener/src/keygen"
	"github.com/asad1123/url-shortener/src/repository/store"
	"github.com/gin-gonic/gin"
)

type App struct {
}

func (a App) Run() {
	config, err := config.LoadConfig(".", "config")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	r := SetupRouter(&config)
	url := fmt.Sprintf("0.0.0.0:%s", config.ServerPort)
	r.Run(url)
}

func SetupRouter(config *config.AppConfig) *gin.Engine {
	handler := startHandler(config)

	r := gin.Default()
	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		// URL manipulation
		urls := v1.Group("/urls")
		urls.POST("", handler.CreateShortenedUrl)
		urls.GET(":id", handler.RetrieveUrlToRedirect)
		urls.DELETE(":id", handler.DeleteShortenedUrl)

		// analytics
		analytics := v1.Group("/analytics/urls")
		analytics.GET(":id", handler.GetUsageAnalyticsForUrl)

		// health check route
		api.GET("/ping", handler.Pong)
	}

	return r
}

func startHandler(config *config.AppConfig) *handlers.Handler {
	store := store.NewStore(config)
	keygen := keygen.NewKeyGen(config)

	return handlers.NewHandler(
		store,
		keygen,
	)
}
