package connector

import (
	"encoding/json"
	"log"

	"github.com/mo0x3f/lark-base-rsshub-sync/pkg/utils"
	"github.com/mo0x3f/lark-base-rsshub-sync/repository/do/connector"
)

const bucketDir = "/root/bucket/tablecache/"

type qcloudRepositoryImpl struct{}

func (repo *qcloudRepositoryImpl) UpdateTable(tableKey string, cache *connector.TableMetaCache) error {
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	err = utils.WriteFile(bucketDir+tableKey, data)
	if err != nil {
		return err
	}
	log.Printf("update table success: %s", cache.String())
	return nil
}

func (repo *qcloudRepositoryImpl) MGetTable(tableKey string) (*connector.TableMetaCache, error) {
	data, err := utils.ReadFile(bucketDir + tableKey)
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
