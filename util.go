package main

import (
	"fmt"
	"log"
	"strings"

	sjson "github.com/bitly/go-simplejson"
)

func SplitGeo(geo string) []string {
	f := func(r rune) bool {
		if r == ' ' || r == ',' {
			return true
		}
		return false
	}
	vals := strings.FieldsFunc(geo, f)
	var ret []string
	for _, val := range vals {
		if val != "" {
			ret = append(ret, val)
		}
	}
	return ret
}

// SearchCityFromLocal ...
func SearchCityFromLocal(search string) (string, error) {
	vals := SplitGeo(search)
	for i, _ := range vals {
		vals[i] = strings.ToLower(vals[i])
	}
	for _, val := range vals {
		for _, city := range CityList {
			if val == city {
				return city, nil
			}
		}
	}
	log.Println("City Not Found: ", search)
	return "", fmt.Errorf("None")
}

// SearchCountryFromLocal ...
func SearchCountryFromLocal(country string) (string, error) {
	vals := SplitGeo(country)
	for _, val := range vals {
		for _, code := range CountryCodeList {
			if val == code {
				return code, nil
			}
		}
	}

	for i, ctry := range CountryList {
		if strings.ToLower(strings.Join(vals, "")) == ctry {
			return CountryCodeList[i], nil
		}
	}
	for name, code := range CountryMap {
		vvals := SplitGeo(name)
		if strings.ToLower(strings.Join(vals, "")) == strings.ToLower(strings.Join(vvals, "")) {
			return code, nil
		}
	}
	log.Println("Country Not Found: ", country)
	return "", fmt.Errorf("None")
}

const (
	API = "https://api.github.com"
)

const (
	ChinaFetchFormat = "/legacy/user/search/location:china?start_page=%d&sort=%s&order=desc"
	WorldFetchFormat = "/legacy/user/search/followers:>0?start_page=%d&sort=%s&order=desc"
)

// china|world
// stars, forks, updated, follower
func SearchUser(location string, page int, attr string) []User {
	var searchURL string
	if location == "china" {
		searchURL = API + fmt.Sprintf(ChinaFetchFormat, page, attr)
	} else {
		searchURL = API + fmt.Sprintf(WorldFetchFormat, page, attr)
	}
	resp, err := client.Get(searchURL)
	if err != nil {
		log.Println("Get error", err)
		return nil
	}
	defer resp.Body.Close()
	js, err := sjson.NewFromReader(resp.Body)
	if err != nil {
		log.Println("Read from body error", err)
		return nil
	}
	log.Println("limit: ", resp.Header.Get("X-RateLimit-Remaining"))

	var users []User
	for i := 0; i < 201; i++ {
		val := js.Get("users").GetIndex(i)
		if IsNilJson(val) {
			break
		}
		u := NewUser(val)
		users = append(users, u)
	}
	return users
}

func NewUser(val *sjson.Json) User {
	user := User{}
	var err error

	user.ID, err = val.Get("id").String()
	if err != nil {
		log.Println("id", err)
	}

	user.Name, err = val.Get("name").String()
	if err != nil {
		log.Println("name", err)
	}

	user.Login, err = val.Get("login").String()
	if err != nil {
		log.Println("login", err)
	}

	user.Location, err = val.Get("location").String()
	if err != nil {
		log.Println("location", err)
	}

	user.FollowersCount, err = val.Get("followers_count").Int()
	if err != nil {
		log.Println("followers_count", err)
	}

	user.Language, err = val.Get("language").String()
	if err != nil {
		log.Println("language", err)
	}

	user.PublicRepoCount, err = val.Get("public_repo_count").Int()
	if err != nil {
		log.Println("public_repo_count", err)
	}

	vals := strings.Split(user.ID, "-")
	if len(vals) >= 2 {
		user.Gravatar = "https://avatars0.githubusercontent.com/u/" + vals[1]
	}
	return user
}

func IsNilJson(val *sjson.Json) bool {
	return val == nil || val.Interface() == nil
}
