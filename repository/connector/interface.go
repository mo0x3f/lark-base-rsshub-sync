package connector

import (
	"fmt"

	"github.com/mo0x3f/lark-base-rsshub-sync/repository/do/connector"
)

type provideRepo func() Repository

type Repository interface {
	UpdateTable(tableKey string, cache *connector.TableMetaCache) error
	MGetTable(tableKey string) (*connector.TableMetaCache, error)
}

var provider provideRepo

func Init(env string) error {
	p, err := getProvider(env)
	if err != nil {
		return err
	}

	provider = p
	return nil
}

func getProvider(env string) (provideRepo, error) {
	switch env {
	case "local":
		return func() Repository {
			return &localCacheRepoImpl{}
		}, nil
	case "replit":
		return func() Repository {
			return &replitRepositoryImpl{}
		}, nil
	case "qcloud":
		return func() Repository {
			return &qcloudRepositoryImpl{}
		}, nil
	}
	return nil, fmt.Errorf("env not found :%s", env)
}

func Get() Repository {
	return provider()
}
