package main

import (
	"testing"
	"time"

	"github.com/google/go-github/github"
)

func TestRateLimit(t *testing.T) {
	limit, _, err := client.RateLimits()
	if err != nil {
		t.Error(err)
	}
	t.Log(limit.String())
	t.Log(client.Rate().String())
}

func TestSearchUsers(t *testing.T) {
	users, resps, err := client.Search.Users("location:china followers:>1000 type:user", &github.SearchOptions{
		Order: "desc",
		Sort:  "followers",
		ListOptions: github.ListOptions{
			Page: 0,
		},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log("remaining: ", resps.Rate.Remaining)
	t.Log("reset: ", resps.Rate.Reset.Sub(time.Now()).Seconds())
	t.Log(len(users.Users))
	t.Log(users.Users[0])
}

func TestCrawlChina(t *testing.T) {
	ret := CrawlChina()
	for user := range ret {
		t.Log(user)
	}
}
