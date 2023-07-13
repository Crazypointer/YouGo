package main

import (
	"YouGo/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig)handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL string `json:"url"`
	}
	jsonDecoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := jsonDecoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Sprintf("Can not decode json: %v", err))
		return
	}
	if params.Name == "" {
		responseWithError(w, http.StatusBadRequest, "The 'name' Can not be empty")
		return
	}
	// 调用数据库接口
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Url: params.URL,
		UserID: user.ID,
	})
	if err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Sprintf("Can not create feed: %v", err))
		return
	}
	responseWithJSON(w, http.StatusOK, databaseFeedToFeed(feed))
}