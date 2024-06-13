package connector

import (
	"encoding/json"
	"strconv"
	"strings"

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

type Request struct {
	Params  *DataSourceConfig `json:"params"`
	Context string            `json:"context"`
}

func (req *Request) UnmarshalJSON(data []byte) error {
	type Alias Request
	alias := &struct {
		*Alias
		Params json.RawMessage `json:"params"`
	}{
		Alias: (*Alias)(req),
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	unquote, err := strconv.Unquote(string(alias.Params))
	if err != nil {
		return err
	}

	// FIXME: 考虑用正则来优化
	unquote = strings.ReplaceAll(unquote, "{\"datasourceConfig\":\"{", "{\"datasourceConfig\":{")
	unquote = strings.ReplaceAll(unquote, "\"}\"}", "\"}}")

	var params struct {
		Config *DataSourceConfig `json:"datasourceConfig"`
	}
	if err = json.Unmarshal([]byte(unquote), &params); err != nil {
		return err
	}

	req.Params = params.Config
	return nil
}

type Response struct {
	Code ResultCode  `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type DataSourceConfig struct {
	RssURL string `json:"rss-url"`
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

func NewFailResponse(code ResultCode, messageKey string) *Response {
	return &Response{
		Code: code,
		Msg:  messageKey,
	}
}

type TableMeta struct {
	TableName string   `json:"tableName"`
	Fields    []*Field `json:"fields"`
}

type Field struct {
	FieldID     string `json:"fieldId"`
	FieldName   string `json:"fieldName"`
	FieldType   int    `json:"fieldType"`
	IsPrimary   bool   `json:"isPrimary"`
	Description string `json:"description"`
	Property    string `json:"property,omitempty"`
}
