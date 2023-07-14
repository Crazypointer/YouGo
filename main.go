package main

import (
	"YouGo/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct{
	DB *database.Queries
}


func main() {
	// 从环境变量中读取配置
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port must be set")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	// 初始化数据库连接
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	db := database.New(conn)
	// 初始化 apiConfig
	apiCfg := apiConfig{
		DB: db, 
	}
	// 初始化路由
	router := chi.NewRouter()
	// 跨域配置
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "OPTIONS", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)

	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows",  apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}",  apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollows))

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	// 添加一个路由组，所有的路由都是以 /v1 为前缀
	router.Mount("/v1", v1Router)
	// 初始化服务
	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	// 启动爬虫
	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(db, collectionConcurrency, collectionInterval)
	log.Printf("Server starting on port %v", portString)
	// 启动并监听服务
	log.Fatal(server.ListenAndServe())
	
}