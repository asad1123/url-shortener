package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUrl(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ping"})
}
