package handler

import (
	"fmt"
	"log"

	"github.com/mo0x3f/lark-base-rsshub-sync/model/connector"
	"github.com/mo0x3f/lark-base-rsshub-sync/pkg/utils"
	"github.com/mo0x3f/lark-base-rsshub-sync/service/rsshub"
)

func (handler *connectorHandlerImpl) ListRecords(req *connector.Request) *connector.Response {
	// 获取订阅链接相关配置信息
	config, err := req.GetValidDataSourceConfig()
	if err != nil {
		log.Println(err.Error())
		return connector.NewFailResponse(connector.ConfigErrCode, connector.ConfigErrorMsg)
	}

	// 请求 RSS 订阅并解析
	log.Println(fmt.Sprintf("target url: %s", config.RssURL))
	feed, err := rsshub.NewService().Fetch(config.RssURL)
	if err != nil {
		log.Println(fmt.Sprintf("rss service err: %s", err.Error()))
		return connector.NewFailResponse(connector.InternalErrCode, connector.InternalErrorMsg)
	}

	// 序列化为 Base 连接器要求的格式
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
			"name": item.Link,
			"url":  item.Link,
		}
		record.Data["author"] = item.Authors
		record.Data["category"] = item.CategoryList
		if item.PubDate != 0 {
			record.Data["pubDate"] = item.PubDate
		}
		result.Records = append(result.Records, record)
	}

	return connector.NewSuccessResponse(result)
}

// 如果 RSS 都没有发布时间，则为全量覆盖模式
func isRSSFeedNoDate(feed *rsshub.Feed) bool {
	// 如果没有拉取到订阅信息，默认为增量
	if len(feed.Items) == 0 {
		return false
	}
	for _, item := range feed.Items {
		if item.PubDate != 0 {
			return false
		}
	}
	return true
}
