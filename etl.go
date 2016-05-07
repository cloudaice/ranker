package main

import (
	"log"
)

type Etler func(<-chan User) <-chan User

func etlCountry(users <-chan User) <-chan User {
	ret := make(chan User)
	go func() {
		defer close(ret)
		for user := range users {
			if u := LoadUser(user.LT, user.ID); u != nil {
				if user.Location == u.Location {
					user.RealCountry = u.RealCountry
				}
			}
			if IsKnownCountryCode(user.RealCountry) {
				ret <- user
				continue
			}
			retC, err := matchCountryFromLocal(user.Location)
			if err == nil {
				user.RealCountry = retC
				ret <- user
				continue
			}
			// log.Printf("match country local error: %s, location: %s\n", err, user.Location)
			retC, err = matchCountryFromInternet(user.Location)
			if err != nil {
				log.Printf("match country internet error: %s, location: %s\n\n", err, user.Location)
				continue
			}
			user.RealCountry = retC
			ret <- user
		}
	}()
	return ret
}

func etlCity(users <-chan User) <-chan User {
	ret := make(chan User)
	go func() {
		defer close(ret)
		for user := range users {
			if user.LT != "china" {
				ret <- user
				continue

			}
			if u := LoadUser(user.LT, user.ID); u != nil {
				if user.Location == u.Location {
					user.RealCity = u.RealCity
				}
			}
			if IsKnownCity(user.RealCity) {
				ret <- user
				continue
			}
			retC, err := matchCityFromLocal(user.Location)
			if err == nil {
				user.RealCity = retC
				ret <- user
				continue
			}
			// log.Printf("match city local error: %s, location: %s\n", err, user.Location)
			retC, err = matchCityFromInternet(user.Location)
			if err != nil {
				log.Printf("match city internet error: %s, location: %s\n", err, user.Location)
				continue
			}
			user.RealCity = retC
			ret <- user
		}

	}()
	return ret
}

func etlScore(users <-chan User) <-chan User {
	ret := make(chan User)
	go func() {
		defer close(ret)
		for user := range users {
			user.Score = user.FollowersCount
			ret <- user
		}
	}()
	return ret
}

func etlLanguage(users <-chan User) <-chan User {
	ret := make(chan User)
	go func() {
		defer close(ret)
		for user := range users {
			u := LoadUser(user.LT, user.ID)
			if u != nil {
				user.Language = u.Language
			}
			ret <- user
		}
	}()
	return ret
}
