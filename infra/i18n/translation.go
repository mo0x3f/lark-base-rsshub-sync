package i18n

import "encoding/json"

var translation = `{
  "err_msg": {
    "internal_error": {
		"en": "Internal error, please contact the plugin owner",
		"zh": "系统异常，插件运行错误"
    },
	"config_error": {
		"en": "invalid params",
      	"zh": "参数异常"
	}
  }
}
`

var translationCache map[string]map[string]map[string]string

func Init() error {
	translationCache = make(map[string]map[string]map[string]string)
	if err := json.Unmarshal([]byte(translation), &translationCache); err != nil {
		return err
	}
	return nil
}

func GetByKey(scope, key string) map[string]string {
	scopeMap := translationCache[scope]
	if scopeMap == nil {
		return nil
	}

	msgMap := scopeMap[key]
	if msgMap == nil {
		return nil
	}

	return msgMap
}
