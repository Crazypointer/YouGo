package main

import (
	"YouGo/internal/database"
	"time"

	"github.com/google/uuid"
)
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string   `json:"name"`
	ApiKey    string   `json:"api_key"`
}


type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string  `json:"name"`
	Url       string `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseUserToUser(dbUser database.User) User{
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:   dbUser.ApiKey,
	}
}
func databaseFeedToFeed(dbFeed database.Feed) Feed{
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		UserID:    dbFeed.UserID,
	}
}

func databaseFeedsToFeeds(dbFeed []database.Feed) []Feed{
	feeds := make([]Feed, 0)
	for _, feed := range dbFeed {
		feeds = append(feeds, databaseFeedToFeed(feed))
	}
	return feeds
}

func databaseFeedFollowToFeedFollow(dbFeedFollw database.FeedFollow) FeedFollow{
	return FeedFollow{
		ID:        dbFeedFollw.ID,
		CreatedAt: dbFeedFollw.CreatedAt,
		UpdatedAt: dbFeedFollw.UpdatedAt,
		UserID:    dbFeedFollw.UserID,
		FeedID:    dbFeedFollw.FeedID,
	}
}
func databaseFeedFollowsToFeedFollows(dbFeedFollw []database.FeedFollow) []FeedFollow{
	feedFollows := make([]FeedFollow, 0)
	for _, feedFollow := range dbFeedFollw {
		feedFollows = append(feedFollows, databaseFeedFollowToFeedFollow(feedFollow))
	}
	return feedFollows
}
type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string   `json:"title"`
	Description *string `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePostsToPosts(dbPosts []database.Post) []Post{
	posts := make([]Post, 0)
	for _, post := range dbPosts {
		var description *string
		if post.Description.Valid {
			description = &post.Description.String
		}
		posts = append(posts, Post{
			ID:          post.ID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			Title:       post.Title,
			Description: description,
			PublishedAt: post.PublishedAt,
			Url:         post.Url,
			FeedID:      post.FeedID,
		})
	}
	return posts
}