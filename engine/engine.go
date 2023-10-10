package engine

import (
	"goweb/common/logx"
	docs "goweb/docs"
	"goweb/engine/middle"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GroupRouterHandler func(*gin.RouterGroup)

func InitEngine(fns ...GroupRouterHandler) *gin.Engine {
	gin.DefaultWriter = logx.Loggerx.Out
	Routerx := gin.New()

	middle.LogFormat(Routerx)
	Routerx.Use(gin.Recovery())
	Routerx.Use(middle.Logging())
	// TODO middleware: i18n, Jwt

	docs.SwaggerInfo.BasePath = "/api/v1"
	Routerx.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := Routerx.Group("/api/v1")
	for _, fn := range fns {
		fn(v1)
	}

	return Routerx
}
