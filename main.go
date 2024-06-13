package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mo0x3f/lark-base-rsshub-sync/handler"
	"github.com/mo0x3f/lark-base-rsshub-sync/middleware"
	"github.com/mo0x3f/lark-base-rsshub-sync/model/connector"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(middleware.VerifySignature())
	r.Use(middleware.RequestLoggerMiddleware())

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"response": "pong",
		})
	})

	// static files
	r.StaticFile("/meta.json", "./assets/meta.json")

	// API: table meta
	r.POST("/api/table_meta", func(c *gin.Context) {
		req := &connector.Request{}
		if err := c.ShouldBind(&req); err != nil {
			log.Println(fmt.Sprintf("parse req err: %+v", err))
			c.JSON(http.StatusOK, connector.NewFailResponse(connector.InternalErrCode, "invalid params"))
			return
		}
		log.Println(fmt.Sprintf("params: %s", req.Params))
		log.Println(fmt.Sprintf("context: %s", req.Context))
		c.JSON(http.StatusOK, handler.NewConnectorHandler().GetTableMeta(req))
	})

	// API: records
	r.POST("/api/records", func(c *gin.Context) {
		c.String(http.StatusOK, "records")
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
