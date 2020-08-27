package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/google/go-github/github"
	"github.com/robfig/cron"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)
	redisService := initRedisClient()
	githubService := initGithubAPIClient()
	cronService := initGronjobClient()

	searchTerm := make(map[string]string)
	searchTerm["machinelearning"] = "python"

	cronService.addGitHubFetchingFunc(redisService, githubService, searchTerm)
	cronService.scheduleAndStart()

	log.Println("Listening on 4000....")
	http.ListenAndServe(":4000", limit(mux))
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

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
