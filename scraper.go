package main

import (
	"YouGo/internal/database"
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)


func startScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration){ 
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequests)
	// 创建一个定时器 定期往一个通道中放时间事件
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C{
		feeds ,err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil{
			log.Println(err)
			continue
		}
		// 同步机制
		wg := &sync.WaitGroup{}
		for _, feed := range feeds{
			wg.Add(1)
			go scrapeFeed(db,wg,feed)
		}		
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed){
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil{
		log.Println("Error marking feed as fetched:",err)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil{
		log.Println("Error fetching feed:", err)
		return
	}
	for _, item := range rssFeed.Channel.Item{
		log.Println("Found post ",item.Title, "on feed", feed.Name)
		description := sql.NullString{}
		if item.Description != ""{
			description.String = item.Description
			description.Valid = true
		}
		// 转换时间(string -> time)
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil{
			log.Println("Error parsing date:", err)
			continue
		}
		_, err = db.CreatePost(context.Background(), 
		database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
    		Description: description,
    		PublishedAt: publishedAt,
    		Url: item.Link,
    		FeedID: feed.ID,
		})
		if err != nil{
			if strings.Contains(err.Error(), "unique constraint"){
				continue
			}
			log.Println("failed to create post:", err)
		}
	}
	log.Printf("Feed %s collected, %v post found", feed.Name, len(rssFeed.Channel.Item))
}