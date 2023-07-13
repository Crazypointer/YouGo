package main

import (
	"YouGo/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// handlerCreateFeedFollow is the handler for the POST /api/feed_follows endpoint.
// It creates a new feed follow for the authenticated user.
// 传入的参数中带有user表明当前handler需要进行权限校验
func (apiCfg *apiConfig)handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	jsonDecoder := json.NewDecoder(r.Body) 
	params := parameters{}
	err := jsonDecoder.Decode(&params)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Sprintf("Can not decode json: %v", err))
		return
	}
	if params.FeedID == uuid.Nil {
		responseWithError(w, http.StatusBadRequest, "The 'feed_id' Can not be empty")
		return
	}
	// 调用数据库接口
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: params.FeedID,
	})
	if err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Sprintf("Can not create feed follow: %v", err))
		return
	}
	responseWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feedFollow))
}
func (apiCfg *apiConfig)handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	
	// 调用数据库接口
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Sprintf("Can not get feed follows: %v", err))
		return
	}
	responseWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig)handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Sprintf("Can not parse feed follow id: %v", err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Sprintf("Can not delete this feed follow: %v", err))
		return
	}
	responseWithJSON(w, http.StatusOK, struct{}{})
}