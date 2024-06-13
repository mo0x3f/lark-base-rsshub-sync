package handler

import "github.com/mo0x3f/lark-base-rsshub-sync/model/connector"

func (handler *connectorHandlerImpl) GetTableMeta(req *connector.Request) *connector.Response {
	meta := &connector.TableMeta{
		TableName: "RSSHub 同步表",
		Fields: []*connector.Field{
			{
				FieldID:     "fld_1",
				FieldName:   "fld_1",
				FieldType:   1,
				IsPrimary:   true,
				Description: "primary fld_1",
			},
			{
				FieldID:     "fld_2",
				FieldName:   "fld_2",
				FieldType:   1,
				IsPrimary:   false,
				Description: "fld_2",
			},
		},
	}
	return connector.NewSuccessResponse(meta)
}
