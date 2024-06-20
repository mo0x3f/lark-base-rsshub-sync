package rsshub

import (
	"fmt"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/mo0x3f/lark-base-rsshub-sync/pkg/utils"
	"github.com/mo0x3f/lark-base-rsshub-sync/service/cache"
)

const defaultCacheExpiration = 10 // 默认缓存时间，单位：分钟

type RSSHubService interface {
	Fetch(subscribeURL string, opts ...WithOptions) (*Feed, error)
}

type WithOptions func(*Config)

func DisableCache(disable bool) WithOptions {
	return func(config *Config) {
		config.DisableCache = disable
	}
}

func CacheExpiration(expiration int) WithOptions {
	return func(config *Config) {
		if expiration > 0 {
			config.CacheExpiration = expiration
		}
	}
}

type Config struct {
	DisableCache    bool // 默认开启缓存
	CacheExpiration int  // 缓存时间
}

type rssHubServiceImpl struct{}

func NewService() RSSHubService {
	return &rssHubServiceImpl{}
}

func (hub *rssHubServiceImpl) Fetch(subscribeURL string, opts ...WithOptions) (*Feed, error) {
	config := &Config{
		DisableCache:    false,
		CacheExpiration: defaultCacheExpiration,
	}
	for _, opt := range opts {
		opt(config)
	}

	fetcher := func() (cache.Serializer, error) {
		parser := gofeed.NewParser()
		result, err := parser.ParseURL(subscribeURL)
		if err != nil {
			return nil, err
		}

		log.Println(fmt.Sprintf("fetch success. Title: %s, Count: %d", result.Title, len(result.Items)))

		feed := &Feed{
			Title: result.Title,
			Items: make([]*Item, 0),
		}
		for _, v := range result.Items {
			item := &Item{}
			item.SetValueWith(v)
			feed.Items = append(feed.Items, item)
		}

		return feed, nil
	}

	if config.DisableCache {
		data, err := fetcher()
		if err != nil {
			return nil, err
		}

		result := data.(*Feed)
		return result, nil
	}

	data := &Feed{}
	err := cache.GetService().GetAndRefresh(utils.Sha256Hash(subscribeURL), data, fetcher, time.Duration(config.CacheExpiration)*time.Minute)
	if err != nil {
		return nil, err
	}
	return data, nil
}
