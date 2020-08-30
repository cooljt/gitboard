package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/google/go-github/github"
	"github.com/robfig/cron"
	"github.com/rs/cors"
)

type redisHandler struct {
	client redisClient
}

func main() {
	mux := http.NewServeMux()

	redisService := initRedisClient()
	githubService := initGithubAPIClient()
	cronService := initGronjobClient()
	redisHandler := redisHandler{
		client: redisService,
	}

	searchTerm := make(map[string]string)
	searchTerm["machinelearning"] = "python"

	cronService.addGitHubFetchingFunc(redisService, githubService, searchTerm)
	cronService.scheduleAndStart()

	mux.HandleFunc("/", redisHandler.pythonHandler)
	log.Println("Listening on 4000....")

	limitHandler := limit(mux)
	finalHanlder := cors.Default().Handler(limitHandler)
	http.ListenAndServe(":4000", finalHanlder)
}

func initRedisClient() redisClient {
	options := &redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}
	client := redisClient{
		client: redis.NewClient(options),
		ctx:    context.Background(),
	}
	return client
}

func initGithubAPIClient() gitHubClient {
	client := gitHubClient{
		client: github.NewClient(nil),
		ctx:    context.Background(),
	}
	return client
}

func initGronjobClient() cronClient {
	client := cronClient{
		client: cron.New(),
	}
	return client
}

func (handler redisHandler) pythonHandler(w http.ResponseWriter, r *http.Request) {
	redisClient := handler.client
	info, _ := redisClient.get("python")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Write([]byte(info))
}
