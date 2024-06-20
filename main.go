package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mo0x3f/lark-base-rsshub-sync/handler"
	"github.com/mo0x3f/lark-base-rsshub-sync/infra/i18n"
	"github.com/mo0x3f/lark-base-rsshub-sync/infra/secret"
	"github.com/mo0x3f/lark-base-rsshub-sync/middleware"
	"github.com/mo0x3f/lark-base-rsshub-sync/model/connector"
	"github.com/mo0x3f/lark-base-rsshub-sync/pkg/flag"
	repo "github.com/mo0x3f/lark-base-rsshub-sync/repository/connector"
	"github.com/mo0x3f/lark-base-rsshub-sync/service/cache"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(middleware.LogLatency())
	r.Use(middleware.RequestLoggerMiddleware())
	r.Use(middleware.VerifySignature())

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

func mustSetupInfra(env string) {
	// 初始化国际化资源
	if err := i18n.Init(); err != nil {
		panic(fmt.Sprintf("i18n.Init() fail: %+v", err))
	}

	// 设置调试环境变量
	flag.SetPageMonitor(os.Getenv("PAGE_MONITOR"))

	// 初始化存储层
	if err := repo.Init(env); err != nil {
		panic(fmt.Sprintf("repo.Init() fail: %+v", err))
	}
}

func mustSetupCache() {
	cacheMode := os.Getenv("CACHE_MODE")
	log.Printf("init with cache mode: %s\n", cacheMode)

	// 设置缓存过期时间，如不设置，则使用默认时间
	flag.SetCacheExpiration(os.Getenv("CACHE_EXPIRATION"))

	if err := cache.Init(cacheMode); err != nil {
		panic(fmt.Sprintf("cache.Init() fail: %+v", err))
	}
}

func mustInitSecret(env string) {
	if err := secret.Init(env); err != nil {
		panic(fmt.Sprintf("secret.Init() fail: %+v", err))
	}
}

func main() {
	// 读取环境变量
	env := os.Getenv("APP_ENV")
	log.Printf("init with env: %s\n", env)

	mustSetupInfra(env)
	mustSetupCache()
	mustInitSecret(env)
	r := setupRouter()
	if err := r.Run(":8080"); err != nil {
		panic(fmt.Sprintf("start server error: %+v", err))
	}
}
