package server

import (
	handlers "github.com/asad1123/url-shortener/server/handlers"
	"github.com/gin-gonic/gin"
)

type Routes struct {
}

func (c Routes) startGin() {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/urls/:id", handlers.GetUrl)
	}
}
