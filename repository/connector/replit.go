package connector

import (
	"encoding/json"
	"log"

	"github.com/mo0x3f/lark-base-rsshub-sync/repository/do/connector"
	"github.com/replit/database-go"
)

type replitRepositoryImpl struct{}

func (repo *replitRepositoryImpl) UpdateTable(tableKey string, cache *connector.TableMetaCache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	err = database.Set(tableKey, string(data))
	if err != nil {
		return err
	}
	log.Printf("update table success: %s", cache.String())
	return nil
}

func (repo *replitRepositoryImpl) MGetTable(tableKey string) (*connector.TableMetaCache, error) {
	data, err := database.Get(tableKey)
	if err != nil {
		return nil, err
	}

	cache := &connector.TableMetaCache{}
	err = json.Unmarshal([]byte(data), &cache)
	if err != nil {
		return nil, err
	}

	return cache, nil
}
