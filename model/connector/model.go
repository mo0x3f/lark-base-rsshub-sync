package connector

type ResultCode int

const (
	SuccessCode       ResultCode = 0
	ConfigErrCode     ResultCode = 1254400
	PermissionErrCode ResultCode = 1254403
	InternalErrCode   ResultCode = 1254500
	PaymentErrCode    ResultCode = 1254505
)

type Request struct {
	Params  string `json:"params"`
	Context string `json:"context"`
}

type Response struct {
	Code ResultCode  `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type DataSourceConfig struct {
	Item1  string `json:"config-item-1"`
	Item2  string `json:"config-item-2"`
	RssURL string `json:"rss-url"`
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
