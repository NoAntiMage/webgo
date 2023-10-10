package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Todo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": 1.0,
	})
}
