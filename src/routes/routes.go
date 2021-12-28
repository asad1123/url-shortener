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
	{
		// URL manipulation
		api.POST("/urls", handler.CreateShortenedUrl)
		api.GET("/urls/:id", handler.RetrieveUrlToRedirect)
		api.DELETE("/urls/:id", handler.DeleteShortenedUrl)

		// analytics
		api.GET("/analytics/urls/:id", handler.GetUsageAnalyticsForUrl)

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
