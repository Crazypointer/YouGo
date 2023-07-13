package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey: 从 http 请求头中获取 API Key
func GetAPIKey(h http.Header) (string, error){
	val := h.Get("Authorization")
	if val == "" {
		return "", errors.New("missing Authorization header")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("invalid Authorization header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("invalid Authorization header")
	}
	return vals[1], nil
}