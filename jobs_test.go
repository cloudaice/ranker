package main

import (
	"testing"
)

func TestRunJobs(t *testing.T) {
	var crl = func() <-chan User {
		ret := make(chan User)
		go func() {
			defer close(ret)
			ret <- User{Score: 0}
			ret <- User{Score: 0}
		}()
		return ret
	}

	var etl = func(users <-chan User) <-chan User {
		ret := make(chan User)
		go func() {
			defer close(ret)
			for user := range users {
				user.Score = user.Score + 1
				ret <- user
			}
		}()
		return ret
	}

	var pst = func(users <-chan User) {
		for user := range users {
			if user.Score != 4 {
				t.Errorf("RunJobs user.Score: %d, exp: %d\n", user.Score, 4)
			}
		}
	}
	RunJobs(crl, pst, etl, etl, etl, etl)
}

func TestRealRunJobs(t *testing.T) {
	var viewUser = func(users <-chan User) {
		for user := range users {
			t.Logf("%v\n", user)
		}
	}
	RunJobs(CrawlChina, viewUser, etlCountry, etlCity, etlScore, etlLanguage)
}
