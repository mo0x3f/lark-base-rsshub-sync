package connector

import (
	"errors"
	"os"

	"github.com/mo0x3f/lark-base-rsshub-sync/repository/do/connector"
)

type Factory interface {
	Init(env string)
	GetRepo() Repository
}

type Repository interface {
	UpdateTable(tableKey string, cache *connector.TableMetaCache) error
	MGetTable(tableKey string) (*connector.TableMetaCache, error)
}

var factory Factory

func Init() error {
	factory = &factoryImpl{}
	factory.Init(os.Getenv("APP_ENV"))
	return nil
}

func GetFactory() Factory {
	return factory
}

var ErrCacheNotExist = errors.New("cache not exist")
