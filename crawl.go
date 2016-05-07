package main

import (
	"log"
	"strconv"
	"time"

	"github.com/google/go-github/github"
)

type Crawler func() <-chan User

func CrawlChina() <-chan User {
	ret := make(chan User)
	query := "location:china followers:>0 type:user"
	go func() {
		defer close(ret)
		var i = 0
		for {
			i = i % 10
			users := searchUser(query, i)
			for _, user := range users {
				user.LT = "china"
				ret <- user
			}
			i++
		}
	}()
	return ret
}

func CrawlWorld() <-chan User {
	ret := make(chan User)
	query := "followers:>0 type:user"
	go func() {
		defer close(ret)
		var i = 0
		for {
			i = i % 10
			users := searchUser(query, i)
			for _, user := range users {
				user.LT = "world"
				ret <- user
			}
			i++
		}
	}()
	return ret
}

func searchUser(query string, page int) (users []User) {
	ret, resp, err := client.Search.Users(query, &github.SearchOptions{
		Order: "desc",
		Sort:  "followers",
		ListOptions: github.ListOptions{
			Page: page,
		},
	})
	if err != nil {
		log.Printf("crawl user error: %s, query: %s, page: %d\n", err, query, page)
		return
	}
	defer resp.Body.Close()
	for _, user := range ret.Users {
		users = append(users, constructUser(user))
	}

	var dela time.Duration = 60 * time.Second
	reset := int(resp.Rate.Reset.Sub(time.Now()).Seconds())
	if reset > 0 {
		dela = time.Duration(resp.Rate.Remaining/reset+5) * time.Second
	}
	log.Printf("Search.Users should sleep %s\n", dela.String())
	time.Sleep(dela)
	return
}

func constructUser(u github.User) (user User) {
	if u.ID != nil {
		user.ID = strconv.Itoa(*u.ID)
	}
	if u.Login != nil {
		user.Login = *u.Login
	}
	client.Users.Get(user.Login)
	realUser, resp, err := client.Users.Get(user.Login)
	if err != nil {
		log.Printf("Get User %s error: %s", user.Login, err)
		return
	}
	defer resp.Body.Close()

	var dela time.Duration = 60 * time.Second
	reset := int(resp.Rate.Reset.Sub(time.Now()).Seconds())
	if reset > 0 {
		dela = time.Duration(resp.Rate.Remaining/reset+5) * time.Second
	}
	log.Printf("Users.Get should sleep %s\n", dela.String())
	time.Sleep(dela)

	if realUser.Name != nil {
		user.Name = *realUser.Name
	}
	if realUser.Location != nil {
		user.Location = *realUser.Location
	}
	if realUser.AvatarURL != nil {
		user.Gravatar = *realUser.AvatarURL
	}
	if realUser.Followers != nil {
		user.FollowersCount = *realUser.Followers
	}
	if realUser.PublicRepos != nil {
		user.PublicRepoCount = *realUser.PublicRepos
	}
	return
}
