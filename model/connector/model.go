package connector

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/mo0x3f/lark-base-rsshub-sync/infra/i18n"
	"github.com/mo0x3f/lark-base-rsshub-sync/pkg/flag"
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

// 默认最大分页数
const defaultMaxPageSize = 1000

type Request struct {
	Params    string         `json:"params"`
	Context   string         `json:"context"`
	ParamsObj *RequestParams `json:"-"`
}

func (req *Request) GetParams() *RequestParams {
	if req.ParamsObj != nil {
		return req.ParamsObj
	}

	obj := &RequestParams{}
	if err := json.Unmarshal([]byte(req.Params), &obj); err != nil {
		log.Printf("request params invalid: %s, %+v\n", req.Params, err)
		return nil
	}

	req.ParamsObj = obj
	return obj
}

func (req *Request) GetValidDataSourceConfig() (*DataSourceConfig, error) {
	params := req.GetParams()
	if params == nil {
		return nil, errors.New("invalid params")
	}

	config := &DataSourceConfig{}
	if err := json.Unmarshal([]byte(params.Config), &config); err != nil {
		return nil, err
	}

	if config != nil && config.Valid() {
		return config, nil
	}

	return nil, errors.New("invalid config")
}

type TableContext struct {
	BaseID   string
	TableID  string
	TenantID string
	UserID   string
}

func GenTableContext(req *Request, url string) (*TableContext, error) {
	if req == nil || url == "" {
		return nil, errors.New("invalid params")
	}

	context := &TableContext{}
	reqCtx := &RequestContext{}
	if err := json.Unmarshal([]byte(req.Context), &reqCtx); err != nil {
		return nil, err
	}

	if reqCtx.ScriptArgs == nil || reqCtx.ScriptArgs.BaseOpenID == "" {
		return nil, errors.New("empty user token")
	}
	context.UserID = reqCtx.ScriptArgs.BaseOpenID

	if reqCtx.TenantKey == "" {
		return nil, errors.New("empty tenant key")
	}
	context.TenantID = reqCtx.TenantKey

	// 一个表只能有一个订阅源，使用订阅源url作为tableID
	context.TableID = utils.Sha256Hash(url)
	return context, nil
}

func (context *TableContext) GetTableKey() string {
	// 租户维度使用 url 做隔离
	return fmt.Sprintf("%s:%s", context.TenantID, context.TableID)
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

func (params *RequestParams) GetNextGUID() string {
	if params.PageToken == "" {
		return ""
	}

	data, err := utils.Base64Decode(params.PageToken)
	if err != nil {
		log.Printf("base64Decode fail. input: %s, %+v\n", params.PageToken, err)
		return ""
	}

	pageToken := &PageToken{}
	if err := json.Unmarshal([]byte(data), &pageToken); err != nil {
		log.Printf("page token invalid: %s, %+v\n", params.PageToken, err)
		return ""
	}

	return pageToken.NextGUID
}

func (params *RequestParams) GetMaxPageSize() int {
	if flag.PageMonitorEnable() {
		return flag.MonitorPageSize()
	}
	if params.MaxPageSize > 0 {
		return params.MaxPageSize
	}
	return defaultMaxPageSize
}

type DataSourceConfig struct {
	RssURL string `json:"rss-url"` // RssURL RSS 订阅链接
}

func (config *DataSourceConfig) Valid() bool {
	return utils.IsValidHTTPURL(config.RssURL)
}

type RequestContext struct {
	Type          string      `json:"type"`
	Bitable       *Bitable    `json:"bitable"`
	ScriptArgs    *ScriptArgs `json:"scriptArgs"`
	PackID        string      `json:"packID"`
	TenantKey     string      `json:"tenantKey"`
	UserTenantKey string      `json:"userTenantKey"`
	BizInstanceID string      `json:"bizInstanceID"`
}

type Bitable struct {
	Token string `json:"token"`
	LogID string `json:"logID"`
}

type ScriptArgs struct {
	ProjectURL string `json:"projectURL"`
	BaseOpenID string `json:"baseOpenID"`
}

type PageToken struct {
	NextGUID string `json:"next_guid"`
}

func GenPageToken(nextGUID string) string {
	token := &PageToken{
		NextGUID: nextGUID,
	}
	data, _ := json.Marshal(token)
	return utils.Base64Encode(data)
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
