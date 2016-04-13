package main

type User struct {
	ID              string `json:"id"`
	Ranker          int    `json: "ranker"`
	Login           string `json:"login"`
	Name            string `json:"name"`
	Location        string `json:"location"`
	Language        string `json:"language"`
	Gravatar        string `json:"gravatar"`
	FollowersCount  int    `json:"followersCount"`
	PublicRepoCount int    `json:"publicRepoCount"`
	Score           int    `json:"score"`
	RealLocation    string `json:"realLocation"`
}

type Location struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type UserList []User

func (ul UserList) Less(i, j int) bool {
	return ul[i].Score > ul[j].Score
}

func (ul UserList) Swap(i, j int) {
	ul[i], ul[j] = ul[j], ul[i]
}

func (ul UserList) Len() int {
	return len(ul)
}
