package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Debug(c *gin.Context) {
	c.String(http.StatusOK, "DEBUG")
}

func Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": 1.0,
	})
}

func Health(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
