package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithError(w http.ResponseWriter, code int, msg string) {
	if code >= 500 {
		log.Println("Response with 5xx error:",msg)
	}
	responseWithJSON(w, code, map[string]interface{}{
		"error": msg,
	})
}


func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// 将 payload 序列化为 JSON 字符串存在一个byte数组中
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		// 如果序列化失败，返回 500
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}
