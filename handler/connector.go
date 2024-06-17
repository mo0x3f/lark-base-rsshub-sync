package handler

import (
	"github.com/mo0x3f/lark-base-rsshub-sync/model/connector"
)

type ConnectorHandler interface {
	GetTableMeta(req *connector.Request) *connector.Response
	ListRecords(req *connector.Request) *connector.Response
}

type connectorHandlerImpl struct{}

func NewConnectorHandler() ConnectorHandler {
	return &connectorHandlerImpl{}
}
