package middleware

import (
	"crypto/sha1"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func VerifySignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// 请求前
		c.Next()

		// 请求后
		latency := time.Since(t)
		log.Print(latency)
	}
}

func GenPostRequestSignature(nonce string, timestamp string, body string, secretKey string) string {
	var b strings.Builder
	b.WriteString(timestamp)
	b.WriteString(nonce)
	b.WriteString(secretKey)
	b.WriteString(body)
	bs := []byte(b.String())
	h := sha1.New()
	h.Write(bs)
	bs = h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
