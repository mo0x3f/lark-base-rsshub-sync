package cache

import (
	"errors"
	"time"
)

type Service interface {
	Init() error
	Get(key string, target Serializer) error
	GetAndRefresh(key string, target Serializer, fetch FetchFunc, expiration time.Duration) error
	Set(key string, value Serializer, expiration time.Duration) error
}

type Serializer interface {
	Serialize() (string, error)
	Deserialize(string) error
	Copy(from Serializer) error
}

type FetchFunc func() (Serializer, error)

func GetService() Service {
	return cacheService
}

var cacheService Service

func Init(mode string) error {
	switch mode {
	case "memory":
		cacheService = &memoryCache{}
		return cacheService.Init()
	default:
		cacheService = &memoryCache{}
		return cacheService.Init()
	}
}

var ErrNotExists = errors.New("cache not found")
