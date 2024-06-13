package handler

import (
	"log"

	"github.com/mmcdole/gofeed"
	"github.com/mo0x3f/lark-base-rsshub-sync/model/connector"
)

func (handler *connectorHandlerImpl) GetTableMeta(req *connector.Request) *connector.Response {
	if !req.Params.Valid() {
		return connector.NewFailResponse(connector.ConfigErrCode, "invalid config")
	}

	rssURL := req.Params.RssURL

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(rssURL)
	if err != nil {
		return connector.NewFailResponse(connector.InternalErrCode, err.Error())
	}
	log.Println(feed.Title)

	meta := &connector.TableMeta{
		TableName: feed.Title,
		Fields: []*connector.Field{
			{
				FieldID:   "title",
				FieldName: "标题",
				FieldType: 1,
				IsPrimary: true,
			},
			{
				FieldID:   "description",
				FieldName: "摘要",
				FieldType: 1,
				IsPrimary: false,
			},
			{
				FieldID:   "link",
				FieldName: "超链接",
				FieldType: 1,
				IsPrimary: false,
			},
			{
				FieldID:   "author",
				FieldName: "作者",
				FieldType: 1,
				IsPrimary: false,
			},
			{
				FieldID:   "pubDate",
				FieldName: "发布时间",
				FieldType: 1,
				IsPrimary: false,
			},
			{
				FieldID:   "category",
				FieldName: "分类",
				FieldType: 1,
				IsPrimary: false,
			},
		},
	}
	return connector.NewSuccessResponse(meta)
}
