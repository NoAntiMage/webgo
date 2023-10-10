package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int
	Message string
	Data    any
}

func FailWithErr(c *gin.Context, statusCode int, err error) {
	resp := Response{}
	resp.Code = 0

	if err != nil {
		resp.Message = err.Error()
	}

	c.JSON(statusCode, resp)
	c.Abort()
}

func SuccessWithData(c *gin.Context, data any) {
	resp := Response{}
	resp.Code = 1
	resp.Message = "ok"
	resp.Data = data

	c.JSON(http.StatusOK, resp)
	c.Abort()
}
