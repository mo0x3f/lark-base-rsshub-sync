package middleware

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		log.Println(fmt.Sprintf("headers: %s", c.Request.Header))
		log.Println(fmt.Sprintf("body: %s", string(body)))
		c.Next()
	}
}
