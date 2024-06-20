package middleware

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mo0x3f/lark-base-rsshub-sync/infra/secret"
	"github.com/mo0x3f/lark-base-rsshub-sync/model/connector"
)

func VerifySignature() gin.HandlerFunc {
	return func(c *gin.Context) {
		signature := c.GetHeader("X-Base-Signature")
		if signature == "" {
			return
		}

		timestamp := c.GetHeader("X-Base-Request-Timestamp")
		nonce := c.GetHeader("X-Base-Request-Nonce")

		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)

		verifySig := genPostRequestSignature(nonce, timestamp, string(body), secret.GetVerifyToken())
		if signature == verifySig {
			c.Next()
			return
		}

		log.Printf("verify signature signature: %s verifySig: %s\n", signature, verifySig)
		c.AbortWithStatusJSON(http.StatusOK, connector.NewFailResponse(connector.PermissionErrCode, connector.VerifyErrorMsg))
	}
}

func genPostRequestSignature(nonce string, timestamp string, body string, secretKey string) string {
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
