package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"net/url"
)

func IsValidHTTPURL(s string) bool {
	parsedURL, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	// 验证是否为http或https协议
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}

func Sha256Hash(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
