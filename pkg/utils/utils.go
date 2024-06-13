package utils

import "net/url"

func IsValidHTTPURL(s string) bool {
	parsedURL, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	// 验证是否为http或https协议
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}
