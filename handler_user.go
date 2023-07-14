package main

import (
	"YouGo/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
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
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
	})
	if err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Sprintf("Can not create user: %v", err))
		return
	}
	responseWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
func (apiCfg *apiConfig)handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	responseWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
func (apiCfg *apiConfig)handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: 10,
	})
	if err != nil {
		responseWithError(w, http.StatusBadRequest, fmt.Sprintf("Can not get posts for user: %v", err))
		return
	}
	responseWithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}