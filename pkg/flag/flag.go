package flag

import (
	"log"
	"strconv"
)

var (
	pageMonitor     = 0 // 模拟单页返回数据，如值为0，则认为以对Lark侧数据为准
	cacheExpiration = 0 // 缓存过期时间（单位：分钟），如值为0，则以模块内默认配置为准
)

func SetPageMonitor(value string) {
	if value == "" {
		return
	}

	pageSize, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("invalid page monitor: %s", value)
	}
	pageMonitor = pageSize
	log.Printf("set page mintior: %d\n", pageMonitor)
}

func PageMonitorEnable() bool {
	return pageMonitor != 0
}

func MonitorPageSize() int {
	return pageMonitor
}

func SetCacheExpiration(value string) {
	if value == "" {
		return
	}
	expiration, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("invalid cache expiration: %s", value)
	}
	cacheExpiration = expiration
	log.Printf("set cache expiration: %d\n", cacheExpiration)
}

func DisableCache() bool {
	return cacheExpiration < 0
}

func GetCacheExpiration() int {
	return cacheExpiration
}
