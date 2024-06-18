package connector

import (
	"fmt"
	"log"
	"sort"
)

type RecordPage []*Record

type TableMetaCache struct {
	ID         string             `json:"table_id"`
	URL        string             `json:"url"`
	RecordMap  map[string]*Record `json:"record_map"`
	RecordPage RecordPage         `json:"-"`
}

func (cache *TableMetaCache) String() string {
	return fmt.Sprintf("[tableID: %s, cache len: %d]", cache.ID, len(cache.RecordMap))
}

func (cache *TableMetaCache) MergeAndSort(items []*Record) bool {
	hasUpdate := false
	for _, item := range items {
		if _, ok := cache.RecordMap[item.Guid]; ok {
			continue
		}
		cache.RecordMap[item.Guid] = item
		log.Printf("add new item: %s\n", item.Guid)
		hasUpdate = true
	}

	// 排序
	cache.RecordPage = cache.SortByTimeASC()

	return hasUpdate
}

// LimitAndSave RecordMap最大值限制，根据时间顺序淘汰
func (cache *TableMetaCache) LimitAndSave(limit int) {
	if len(cache.RecordPage) <= limit {
		return
	}

	sort.Slice(cache.RecordPage, func(i, j int) bool {
		return cache.RecordPage[i].PubDate > cache.RecordPage[j].PubDate
	})

	log.Printf("origin: %d; cut: %d\n", len(cache.RecordPage), len(cache.RecordPage)-limit)

	// cut to limit
	limitPage := cache.RecordPage[:limit]

	// delete map value
	newRecordMap := make(map[string]*Record)
	for _, item := range limitPage {
		newRecordMap[item.Guid] = item
	}
	cache.RecordMap = newRecordMap
}

func (cache *TableMetaCache) SortByTimeASC() RecordPage {
	collection := make(RecordPage, 0)
	for _, item := range cache.RecordMap {
		collection = append(collection, item)
	}

	sort.Slice(collection, func(i, j int) bool {
		return collection[i].PubDate < collection[j].PubDate
	})

	return collection
}

func (collection RecordPage) IndexOfGuid(guid string) int {
	for idx, record := range collection {
		if record.Guid == guid {
			return idx
		}
	}
	return -1
}

func (collection RecordPage) NextPage(guid string, pageSize int) ([]*Record, string) {
	if len(collection) <= pageSize {
		return collection, ""
	}

	// 返回第一页数据
	if guid == "" {
		return collection[:pageSize], collection[pageSize].Guid
	}

	// 获取排序后的 Guid 的索引
	index := collection.IndexOfGuid(guid)
	if index == -1 {
		log.Printf("fetch first page. guid index not found: %s\n", guid)
		return collection[:pageSize], collection[pageSize].Guid
	}

	// 分页返回数据
	if index+pageSize >= len(collection) {
		return collection[index:], ""
	} else {
		return collection[index : index+pageSize], collection[index+pageSize].Guid
	}
}

type Record struct {
	Guid         string   `json:"guid"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Link         string   `json:"link"`
	Authors      []string `json:"authors"`
	PubDate      int64    `json:"pubDate"`
	CategoryList []string `json:"category_list"`
}
