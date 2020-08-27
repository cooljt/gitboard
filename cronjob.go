package main

import (
	"encoding/json"
	"log"

	"github.com/robfig/cron"
)

type cronClient struct {
	client *cron.Cron
}

func (c cronClient) addGitHubFetchingFunc(redisClient redisClient, githubClient gitHubClient, searchTerm map[string]string) {
	c.client.AddFunc("@every 1h", func() {
		githubFetchAndStore(redisClient, githubClient, searchTerm)
	})
}

func (c cronClient) scheduleAndStart() {
	c.client.Start()
}

func githubFetchAndStore(redisClient redisClient, githubClient gitHubClient, searchTerm map[string]string) {
	for k, v := range searchTerm {
		repos := githubClient.searchGithubRepoByStar(10, k, v)
		log.Println("Get Github Result...")
		if repos != nil {

			reposStr, err := json.Marshal(repos)
			if err != nil {
				continue
			}

			err = redisClient.set(v, string(reposStr))
			if err != nil {
				continue
			}
			log.Println("Save record...")
		}
	}
}
