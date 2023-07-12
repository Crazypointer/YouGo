package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	responseJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}
