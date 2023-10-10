package router

import (
	"goweb/app/info/api"

	"github.com/gin-gonic/gin"
)

func InitInfoRouter(Router *gin.RouterGroup) {
	infoRouter := Router.Group("info")
	{
		infoRouter.GET("/debug", api.Debug)
		infoRouter.GET("/version", api.Version)
		infoRouter.GET("/health", api.Health)
	}
}
