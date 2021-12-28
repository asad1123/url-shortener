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
	r := gin.Default()

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	handler := startHandler(&config)

	api := r.Group("/api")
	{
		api.POST("/urls", handler.CreateShortenedUrl)
		api.GET("/urls/:id", handler.RetrieveUrlToRedirect)
		api.DELETE("/urls/:id", handler.DeleteShortenedUrl)

		api.GET("/analytics/urls/:id", handler.GetUsageAnalyticsForUrl)
	}

	url := fmt.Sprintf("0.0.0.0:%s", config.ServerPort)
	r.Run(url)
}

func startHandler(config *config.AppConfig) *handlers.Handler {
	store := store.NewStore(config)
	keygen := keygen.NewKeyGen(config)

	return handlers.NewHandler(
		store,
		keygen,
	)
}
