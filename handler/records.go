package handler

import (
	"fmt"
	"log"

	"github.com/mo0x3f/lark-base-rsshub-sync/model/connector"
	"github.com/mo0x3f/lark-base-rsshub-sync/pkg/flag"
	"github.com/mo0x3f/lark-base-rsshub-sync/pkg/utils"
	repo "github.com/mo0x3f/lark-base-rsshub-sync/repository/connector"
	do "github.com/mo0x3f/lark-base-rsshub-sync/repository/do/connector"
	"github.com/mo0x3f/lark-base-rsshub-sync/service/rsshub"
)

const maxRecordLimit = 10000

func (handler *connectorHandlerImpl) ListRecords(req *connector.Request) *connector.Response {
	// 获取订阅链接相关配置信息
	config, err := req.GetValidDataSourceConfig()
	if err != nil {
		log.Println(err.Error())
		return connector.NewFailResponse(connector.ConfigErrCode, connector.ConfigErrorMsg)
	}

	// 生成 table 上下文
	tableCtx, err := connector.GenTableContext(req, config.RssURL)
	if err != nil {
		log.Println(err.Error())
		return connector.NewFailResponse(connector.ConfigErrCode, connector.InternalErrorMsg)
	}

	// 请求 RSS 订阅并解析
	feed, err := rsshub.NewService().Fetch(config.RssURL, rsshub.DisableCache(flag.DisableCache()), rsshub.CacheExpiration(flag.GetCacheExpiration()))
	if err != nil {
		log.Println(fmt.Sprintf("rss service err: %s", err.Error()))
		return connector.NewFailResponse(connector.InternalErrCode, connector.InternalErrorMsg)
	}

	if len(feed.Items) == 0 {
		log.Println("no feed items")
		return connector.NewFailResponse(connector.InternalErrCode, connector.InternalErrorMsg)
	}

	// 如果订阅内容没有发布时间，则认为全量覆盖，不需要做分页与存储
	overrideMode := feed.IsOverrideMode()
	log.Printf("override mode: %v\n", overrideMode)

	if overrideMode {
		result := &connector.RecordsPage{
			HasMore: false,
			Records: feed2RecordList(feed),
		}
		return connector.NewSuccessResponse(result)
	}

	// 查询之前缓存的 Feed
	tableKey := tableCtx.GetTableKey()
	log.Printf("route tableKey: %s\n", tableKey)

	tableCache, err := repo.Get().MGetTable(tableCtx.GetTableKey())
	if err != nil {
		// 查询缓存失败，不报错，继续执行
		log.Printf("MGetTable fail err: %+v\n", err)
	}

	// 第一次请求
	if tableCache == nil {
		tableCache = newCacheWithFeed(tableKey, config.RssURL, feed)
		err = repo.Get().UpdateTable(tableKey, tableCache)
		if err != nil {
			// 更新缓存失败，不报错，继续返回订阅查询结果
			log.Printf("update table err: %s, %+v\n", tableKey, err)
		}

		result := &connector.RecordsPage{
			HasMore: false,
			Records: feed2RecordList(feed),
		}
		return connector.NewSuccessResponse(result)
	}

	// 有缓存，先合并数据
	recordDOs := make([]*do.Record, 0)
	for _, item := range feed.Items {
		recordDOs = append(recordDOs, feedItem2RecordDO(item))
	}
	// 合并 & 按照 Feed 发布时间顺序排序
	hasUpdate := tableCache.MergeAndSort(recordDOs)

	// 更新缓存
	if hasUpdate {
		tableCache.LimitAndSave(maxRecordLimit)
		if err = repo.Get().UpdateTable(tableKey, tableCache); err != nil {
			// 更新缓存失败，不报错，继续返回订阅查询结果
			log.Printf("merge and update table err: %s, %+v\n", tableKey, err)
		}
	}

	// 分页返回数据
	guid := req.GetParams().GetNextGUID()
	log.Printf("next guid: %s\n", guid)
	perPage, nextGuid := tableCache.RecordPage.NextPage(guid, req.GetParams().GetMaxPageSize())

	result := &connector.RecordsPage{
		HasMore:       nextGuid != "", // 如果返回了nextGuid，则认为还有下一页
		NextPageToken: connector.GenPageToken(nextGuid),
		Records:       recordDO2RecordList(perPage),
	}
	return connector.NewSuccessResponse(result)
}

func feed2RecordList(feed *rsshub.Feed) []*connector.Record {
	records := make([]*connector.Record, 0)
	for _, item := range feed.Items {
		record := feedItem2Record(item)
		records = append(records, record)
	}
	return records
}

func feedItem2Record(item *rsshub.Item) *connector.Record {
	record := &connector.Record{
		PrimaryID: utils.Sha256Hash(item.GUID),
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
	return record
}

func feedItem2RecordDO(item *rsshub.Item) *do.Record {
	return &do.Record{
		Guid:         item.GUID,
		Title:        item.Title,
		Description:  item.Description,
		Link:         item.Link,
		Authors:      item.Authors,
		PubDate:      item.PubDate,
		CategoryList: item.CategoryList,
	}
}

func recordDO2RecordList(items []*do.Record) []*connector.Record {
	records := make([]*connector.Record, 0)
	for _, item := range items {
		record := recordDO2Record(item)
		records = append(records, record)
	}
	return records
}

func recordDO2Record(item *do.Record) *connector.Record {
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
	return record
}

func newCacheWithFeed(tableID string, url string, feed *rsshub.Feed) *do.TableMetaCache {
	cache := &do.TableMetaCache{
		ID:        tableID,
		URL:       url,
		RecordMap: make(map[string]*do.Record),
	}

	for _, item := range feed.Items {
		cache.RecordMap[item.GUID] = feedItem2RecordDO(item)
	}

	return cache
}
