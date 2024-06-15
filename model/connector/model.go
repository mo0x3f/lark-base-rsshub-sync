package connector

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/mo0x3f/lark-base-rsshub-sync/infra/i18n"
	"github.com/mo0x3f/lark-base-rsshub-sync/pkg/utils"
)

type ResultCode int

const (
	SuccessCode       ResultCode = 0
	ConfigErrCode     ResultCode = 1254400
	PermissionErrCode ResultCode = 1254403
	InternalErrCode   ResultCode = 1254500
	PaymentErrCode    ResultCode = 1254505
)

type MessageKey string

const (
	InternalErrorMsg MessageKey = "internal_error"
	ConfigErrorMsg   MessageKey = "config_error"
)

// 兜底错误提示
const defaultErrorMsg = "{\"en\":\"Internal error\",\"zh\":\"系统异常，插件运行错误\"}"

type Request struct {
	Params  string `json:"params"`
	Context string `json:"context"`
}

func (req *Request) GetValidDataSourceConfig() (*DataSourceConfig, error) {
	var params = &struct {
		ConfigStr string `json:"datasourceConfig"`
	}{}

	if err := json.Unmarshal([]byte(req.Params), &params); err != nil {
		return nil, err
	}

	config := &DataSourceConfig{}
	if err := json.Unmarshal([]byte(params.ConfigStr), &config); err != nil {
		return nil, err
	}

	if config != nil && config.Valid() {
		return config, nil
	}

	return nil, errors.New("invalid config")
}

type Response struct {
	Code ResultCode  `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type RequestParams struct {
	Config        string `json:"datasourceConfig"`
	TransactionID string `json:"transactionID"`
	PageToken     string `json:"pageToken"`
	MaxPageSize   int    `json:"maxPageSize"`
}

type DataSourceConfig struct {
	RssURL string `json:"rss-url"` // RSS 订阅链接
}

func (config *DataSourceConfig) Valid() bool {
	return utils.IsValidHTTPURL(config.RssURL)
}

type RequestContext struct {
	Bitable       *Bitable `json:"bitable"`
	ScriptArgs    *Bitable `json:"scriptArgs"`
	PackID        string   `json:"packID"`
	TenantKey     string   `json:"tenantKey"`
	UserTenantKey string   `json:"userTenantKey"`
}

type Bitable struct {
	Token string `json:"token"`
	LogID string `json:"logID"`
}

type ScriptArgs struct {
	ProjectURL string `json:"projectURL"`
	BaseOpenID string `json:"baseOpenID"`
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Code: SuccessCode,
		Msg:  "",
		Data: data,
	}
}

func NewFailResponse(code ResultCode, messageKey MessageKey) *Response {
	msg := i18n.GetByKey("err_msg", string(messageKey))
	if msg == nil {
		log.Printf("message not found: %s\n", messageKey)
		return &Response{
			Code: code,
			Msg:  defaultErrorMsg,
		}
	}

	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Printf("json marshal fail: %s\n", err.Error())
		return &Response{
			Code: code,
			Msg:  defaultErrorMsg,
		}
	}

	return &Response{
		Code: code,
		Msg:  string(msgByte),
	}
}

type TableMeta struct {
	TableName string   `json:"tableName"`
	Fields    []*Field `json:"fields"`
}

type Field struct {
	FieldID     string    `json:"fieldId"`
	FieldName   string    `json:"fieldName"`
	FieldType   int       `json:"fieldType"`
	IsPrimary   bool      `json:"isPrimary"`
	Description string    `json:"description"`
	Property    *Property `json:"property,omitempty"`
}

type Property struct {
	Formatter string `json:"formatter,omitempty"`
}

type RecordsPage struct {
	NextPageToken string    `json:"nextPageToken"`
	HasMore       bool      `json:"hasMore"`
	Records       []*Record `json:"records"`
}

type Record struct {
	PrimaryID string                 `json:"primaryID"`
	Data      map[string]interface{} `json:"data"`
}
