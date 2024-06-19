package rsshub

import (
	"encoding/json"
	"errors"

	"github.com/mmcdole/gofeed"
	"github.com/mo0x3f/lark-base-rsshub-sync/service/cache"
)

type Feed struct {
	Title string  `json:"title"`
	Items []*Item `json:"items"`
}

func (feed *Feed) Serialize() (string, error) {
	data, err := json.Marshal(feed)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (feed *Feed) Deserialize(data string) error {
	item := &Feed{}
	if err := json.Unmarshal([]byte(data), &item); err != nil {
		return err
	}
	feed.Title = item.Title
	feed.Items = item.Items
	return nil
}

func (feed *Feed) Copy(from cache.Serializer) error {
	obj, ok := from.(*Feed)
	if !ok {
		return errors.New("types error")
	}
	feed.Title = obj.Title
	feed.Items = obj.Items
	return nil
}

// IsOverrideMode 如果 RSS 都没有发布时间，则为全量覆盖模式
func (feed *Feed) IsOverrideMode() bool {
	for _, item := range feed.Items {
		if item.PubDate != 0 {
			return false
		}
	}
	return true
}

type Item struct {
	GUID         string   `json:"guid"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Link         string   `json:"link"`
	Authors      []string `json:"authors"`
	PubDate      int64    `json:"pub_date"`
	CategoryList []string `json:"category_list"`
}

func (item *Item) SetValueWith(v *gofeed.Item) {
	item.GUID = v.GUID
	item.Title = v.Title
	item.Description = v.Description
	item.Link = v.Link

	if v.PublishedParsed == nil {
		item.PubDate = 0
	} else {
		item.PubDate = v.PublishedParsed.UnixMilli()
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
}
