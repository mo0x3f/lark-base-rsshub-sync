package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mo0x3f/lark-base-rsshub-sync/handler"
	"github.com/mo0x3f/lark-base-rsshub-sync/infra/i18n"
	"github.com/mo0x3f/lark-base-rsshub-sync/middleware"
	"github.com/mo0x3f/lark-base-rsshub-sync/model/connector"
	"github.com/mo0x3f/lark-base-rsshub-sync/pkg/flag"
	repo "github.com/mo0x3f/lark-base-rsshub-sync/repository/connector"
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
		if err := c.ShouldBind(req); err != nil {
			log.Println(fmt.Sprintf("parse req err: %+v", err))
			c.JSON(http.StatusOK, connector.NewFailResponse(connector.ConfigErrCode, connector.ConfigErrorMsg))
			return
		}
		c.JSON(http.StatusOK, handler.NewConnectorHandler().GetTableMeta(req))
	})

	// API: records
	r.POST("/api/records", func(c *gin.Context) {
		req := &connector.Request{}
		if err := c.ShouldBind(req); err != nil {
			log.Println(fmt.Sprintf("parse req err: %+v", err))
			c.JSON(http.StatusOK, connector.NewFailResponse(connector.ConfigErrCode, connector.ConfigErrorMsg))
			return
		}
		c.JSON(http.StatusOK, handler.NewConnectorHandler().ListRecords(req))
	})

	return r
}

func mustSetupInfra() {
	// 初始化国际化资源
	if err := i18n.Init(); err != nil {
		panic(fmt.Sprintf("i18n.Init() fail: %+v", err))
	}

	// 读取环境变量
	env := os.Getenv("APP_ENV")
	log.Printf("init with env: %s\n", env)

	// 设置调试环境变量
	flag.SetPageMonitor(os.Getenv("PAGE_MONITOR"))

	// 初始化存储层
	if err := repo.Init(env); err != nil {
		panic(fmt.Sprintf("repo.Init() fail: %+v", err))
	}
}

func main() {
	mustSetupInfra()
	r := setupRouter()
	r.Run(":8080")
}
