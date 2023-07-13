package main

import (
	"YouGo/internal/auth"
	"YouGo/internal/database"
	"fmt"
	"net/http"
)

type authheadHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig)middlewareAuth(handler authheadHandler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中获取apiKey
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			responseWithError(w, http.StatusForbidden, fmt.Sprintf("Auth error: %v", err))
			return
		}
		// 调用数据库接口 根据apiKey获取用户信息
		user, err := apiCfg.DB.GetUser(r.Context(), apiKey)
		if err != nil {
			responseWithError(w, http.StatusBadRequest, fmt.Sprintf("User not found: %v", err))
			return
		}
		// 调用传入的handler
		handler(w, r, user)
	}
}