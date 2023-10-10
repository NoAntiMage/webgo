package middle

import (
	"bytes"
	"fmt"
	"goweb/common/logx"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var webSuffix []string = []string{"html", "css", "js", "png"}

type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.b.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogFormat(engine *gin.Engine) {
	engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s web http:\"%s %s %d %s %s\"\n",
			param.TimeStamp.Format(time.DateTime),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ErrorMessage,
		)
	}))
}

//as a restful api server, logging request & respone.
//ignore web suffix. e.g: swaggerUI element.
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ignoreFlag bool

		respWriter := responseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}

		for _, suffix := range webSuffix {
			if strings.HasSuffix(c.Request.URL.Path, suffix) {
				ignoreFlag = true
			}
		}

		if ignoreFlag != true {
			if c.Request.Method == "POST" {
				reqBody, err := c.GetRawData()
				if err == nil {
					logx.Loggerx.Infof("Request -> %s", string(reqBody))
				} else {
					logx.Loggerx.Infof("Request ERR %v", err)
				}
				c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			c.Writer = respWriter
		}

		c.Next()

		if ignoreFlag != true {
			logx.Loggerx.Infof("Response <- %s", respWriter.b.String())
		}
	}
}
