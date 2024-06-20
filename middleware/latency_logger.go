package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LogLatency() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// 请求前
		c.Next()

		// 请求后
		latency := time.Since(t)
		log.Printf("%s: %v\n", c.Request.URL, latency)
	}
}
