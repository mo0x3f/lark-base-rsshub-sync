package handler

import "github.com/mo0x3f/lark-base-rsshub-sync/model/connector"

func (handler *connectorHandlerImpl) ListRecords(req *connector.Request) *connector.Response {
	return connector.NewSuccessResponse(nil)
}
