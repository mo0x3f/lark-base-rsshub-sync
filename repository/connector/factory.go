package connector

import (
	"fmt"
)

type factoryImpl struct {
	Env string
}

func (factory *factoryImpl) Init(env string) {
	factory.Env = env
}

func (factory *factoryImpl) GetRepo() Repository {
	switch factory.Env {
	case "local":
		return &localCacheRepoImpl{}
	case "replit":
		return &replitRepositoryImpl{}
	}
	panic(fmt.Sprintf("env not found: %s", factory.Env))
}
