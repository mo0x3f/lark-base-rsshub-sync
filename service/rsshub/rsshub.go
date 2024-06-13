package rsshub

import (
	"fmt"
	"log"

	"github.com/mmcdole/gofeed"
)

type RSSHubService interface {
	Fetch(subscribeURL string) (*Feed, error)
}

type Feed struct {
	Title string
	Items []*Item
}

type Item struct {
	Guid         string
	Title        string
	Description  string
	Link         string
	Authors      []string
	PubDate      int64
	CategoryList []string
}

type rssHubServiceImpl struct{}

func NewService() RSSHubService {
	return &rssHubServiceImpl{}
}

func (hub *rssHubServiceImpl) Fetch(subscribeURL string) (*Feed, error) {
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
		item := &Item{
			Guid:        v.GUID,
			Title:       v.Title,
			Description: v.Description,
			Link:        v.Link,
			PubDate:     v.PublishedParsed.UnixMicro(),
		}

		for _, author := range v.Authors {
			if len(item.Authors) == 0 {
				item.Authors = make([]string, 0)
			}
			item.Authors = append(item.Authors, author.Name)
		}

		for _, category := range v.Categories {
			if len(item.CategoryList) == 0 {
				item.CategoryList = make([]string, 0)
			}
			item.CategoryList = append(item.CategoryList, category)
		}

		feed.Items = append(feed.Items, item)
	}

	return feed, nil
}
