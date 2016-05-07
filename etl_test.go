package main

import (
	"testing"
)

func TestEtl(t *testing.T) {
	source := make(chan User)
	ret := etlCity(etlCountry(source))
	source <- User{
		ID:              "290496",
		Ranker:          0,
		Login:           "lepture",
		Name:            "Hsiaoming Yang",
		Location:        "beijing, China",
		Language:        "JavaScript",
		Gravatar:        "xxx.com/me.jpg",
		FollowersCount:  1000,
		PublicRepoCount: 10,
		Score:           10,
		RealCity:        "",
		RealCountry:     "",
	}
	close(source)
	u := <-ret
	if u.RealCity != "beijing" {
		t.Error(u.RealCity)
	}
	if u.RealCountry != "CN" {
		t.Error(u.RealCountry)
	}
}

func TestEtlReal(t *testing.T) {
	ret := etlCity(etlCountry(CrawlChina()))
	for user := range ret {
		t.Log(user)
	}
}
