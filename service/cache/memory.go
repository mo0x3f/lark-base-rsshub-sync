package cache

import (
	"fmt"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

type memoryCache struct {
	cache *cache.Cache
}

func (c *memoryCache) Init() error {
	c.cache = cache.New(10*time.Minute, 15*time.Minute)
	return nil
}

func (c *memoryCache) Get(key string, target Serializer) error {
	val, found := c.cache.Get(key)
	if !found {
		return errors.Wrap(ErrNotExists, fmt.Sprintf("not found: %s", key))
	}

	data, ok := val.(string)
	if !ok {
		return fmt.Errorf("type err: %s", key)
	}

	log.Printf("cache hit: %s\n", key)

	err := target.Deserialize(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *memoryCache) GetAndRefresh(key string, target Serializer, fetch FetchFunc, expiration time.Duration) error {
	err := c.Get(key, target)
	if err == nil {
		return nil
	}

	if fetch == nil {
		return err
	}

	value, err := fetch()
	if err != nil {
		return err
	}

	err = c.Set(key, value, expiration)
	if err != nil {
		fmt.Printf("save cache err: %s: %+v\n", key, err)
	}

	err = target.Copy(value)
	if err != nil {
		return err
	}
	return nil
}

func (c *memoryCache) Set(key string, value Serializer, expiration time.Duration) error {
	data, err := value.Serialize()
	if err != nil {
		return err
	}

	c.cache.Set(key, data, expiration)
	return nil
}
