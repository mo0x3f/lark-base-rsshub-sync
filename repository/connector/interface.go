package connector

import (
	"errors"

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

func Init(env string) error {
	factory = &factoryImpl{}
	factory.Init(env)
	return nil
}

func GetFactory() Factory {
	return factory
}

var ErrCacheNotExist = errors.New("cache not exist")
