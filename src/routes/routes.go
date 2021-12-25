package routes

import (
	handlers "github.com/asad1123/url-shortener/src/handlers"
	"github.com/gin-gonic/gin"
)

type App struct {
}

func (a App) Run() {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/urls", handlers.CreateShortenedUrl)
		api.GET("/urls/:id", handlers.RetrieveShortenedUrl)
		api.DELETE("/urls/:id", handlers.DeleteShortenedUrl)

		api.GET("/analytics/urls/:id", handlers.GetUsageAnalyticsForUrl)
	}
	r.Run(":8000")
}
