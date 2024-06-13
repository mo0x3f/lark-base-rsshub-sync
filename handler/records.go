package handler

import (
	"fmt"
	"log"

	"github.com/mo0x3f/lark-base-rsshub-sync/model/connector"
	"github.com/mo0x3f/lark-base-rsshub-sync/pkg/utils"
	"github.com/mo0x3f/lark-base-rsshub-sync/service/rsshub"
)

func (handler *connectorHandlerImpl) ListRecords(req *connector.Request) *connector.Response {
	config, err := req.GetValidDataSourceConfig()
	if err != nil {
		log.Println(err.Error())
		return connector.NewFailResponse(connector.ConfigErrCode, "invalid config")
	}

	log.Println(fmt.Sprintf("target url: %s", config.RssURL))

	feed, err := rsshub.NewService().Fetch(config.RssURL)
	if err != nil {
		log.Println(fmt.Sprintf("rss service err: %s", err.Error()))
		return connector.NewFailResponse(connector.InternalErrCode, err.Error())
	}

	result := &connector.RecordsPage{
		HasMore: false,
		Records: make([]*connector.Record, 0),
	}
	for _, item := range feed.Items {
		record := &connector.Record{
			PrimaryID: utils.Sha256Hash(item.Guid),
			Data:      make(map[string]interface{}),
		}
		record.Data["title"] = item.Title
		record.Data["description"] = item.Description
		record.Data["link"] = map[string]string{
			"name": "跳转链接",
			"url":  item.Link,
		}
		record.Data["author"] = item.Authors
		record.Data["pubDate"] = item.PubDate
		record.Data["category"] = item.CategoryList
		result.Records = append(result.Records, record)
	}

	return connector.NewSuccessResponse(result)
}
