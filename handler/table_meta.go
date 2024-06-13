package handler

import (
	"fmt"
	"log"

	"github.com/mo0x3f/lark-base-rsshub-sync/model/connector"
	"github.com/mo0x3f/lark-base-rsshub-sync/service/rsshub"
)

func (handler *connectorHandlerImpl) GetTableMeta(req *connector.Request) *connector.Response {
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
				FieldType: 10,
				IsPrimary: false,
			},
			{
				FieldID:   "author",
				FieldName: "作者",
				FieldType: 4,
				IsPrimary: false,
			},
			{
				FieldID:   "pubDate",
				FieldName: "发布时间",
				FieldType: 5,
				IsPrimary: false,
				Property: &connector.Property{
					Formatter: "yyyy/MM/dd",
				},
			},
			{
				FieldID:   "category",
				FieldName: "分类",
				FieldType: 4,
				IsPrimary: false,
			},
		},
	}
	return connector.NewSuccessResponse(meta)
}
