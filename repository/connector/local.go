package connector

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/mo0x3f/lark-base-rsshub-sync/repository/do/connector"
)

var tableCache = map[string]string{}

type localCacheRepoImpl struct{}

func (repo *localCacheRepoImpl) UpdateTable(tableKey string, cache *connector.TableMetaCache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	tableCache[tableKey] = string(data)
	log.Printf("update table success: %s", cache.String())
	return nil
}

func (repo *localCacheRepoImpl) MGetTable(tableKey string) (*connector.TableMetaCache, error) {
	val, ok := tableCache[tableKey]
	if !ok {
		return nil, errors.New("not found")
	}

	table := &connector.TableMetaCache{}
	if err := json.Unmarshal([]byte(val), &table); err != nil {
		return nil, err
	}

	return table, nil
}
