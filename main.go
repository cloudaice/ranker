package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
)

func main() {
	//fmt.Println(MatchCity("taizhou linhai"))
	//fmt.Println(SearchUser("china", 1, "followers"))
	//fmt.Println(SearchUser("world", 1, "followers"))
	//RunServer()
	//	source := &UserSearcher{}
	//	worker := &UserSyncer{origin: source}
	//	worker.Start()
	//log.Println(source.GetChinaUser())
	//log.Println(source.GetWorldUser())
	RunServer()
	//StoreWorker()
}

func RunServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/*page", IndexHandler)
	mux.HandleFunc("/", IndexHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	mux.HandleFunc("/githubchina", ChinaHandler)
	mux.HandleFunc("/githubworld", WorldHandler)
	mux.HandleFunc("/chinamap", ChinaMapHandler)
	mux.HandleFunc("/worldmap", WorldMapHandler)
	mux.HandleFunc("/favicon.ico", FaviconHandler)
	log.Fatal(http.ListenAndServe(":9090", mux))
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
}

func Tpl(name string) string {
	return "./template/" + name
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadFile(Tpl("index.html"))
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Parse tmpls error: %s", err)))
		return
	}
	w.Write(body)
}

func ParsePage(r *http.Request) (int, int) {
	currentPage, err := strconv.ParseInt(r.FormValue("current_page"), 10, 64)
	if err != nil {
		currentPage = 1
	}
	pageSize, err := strconv.ParseInt(r.FormValue("page_size"), 10, 64)
	if err != nil {
		pageSize = 10
	}
	return int(currentPage), int(pageSize)
}

func LoadUser(data [][]byte) []User {
	var users []User
	for _, val := range data {
		user := User{}
		if err := json.Unmarshal(val, &user); err != nil {
			log.Println("Unmarshal error", err)
		}
		users = append(users, user)
	}
	return users
}

func ChinaHandler(w http.ResponseWriter, r *http.Request) {
	currentPage, pageSize := ParsePage(r)
	data, err := LoadBucket("china")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	users := LoadUser(data)
	for i, _ := range users {
		users[i].Score = users[i].FollowersCount
	}
	sort.Sort(UserList(users))
	for i, _ := range users {
		users[i].Ranker = i + 1
	}
	if currentPage*pageSize > len(users) {
		pageSize = 0
	}

	body, err := json.MarshalIndent(map[string]interface{}{
		"status": 0,
		"total":  len(users),
		"users":  users[(currentPage-1)*pageSize : currentPage*pageSize],
	}, "", "    ")
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(body)
}

func WorldHandler(w http.ResponseWriter, r *http.Request) {
	currentPage, pageSize := ParsePage(r)
	data, err := LoadBucket("world")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	users := LoadUser(data)
	for i, _ := range users {
		users[i].Score = users[i].FollowersCount
	}

	sort.Sort(UserList(users))
	for i, _ := range users {
		users[i].Ranker = i + 1
	}
	if currentPage*pageSize > len(users) {
		pageSize = 0
	}
	body, err := json.MarshalIndent(map[string]interface{}{
		"status": 0,
		"total":  len(users),
		"users":  users[(currentPage-1)*pageSize : currentPage*pageSize],
	}, "", "    ")
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(body)

}

func ChinaMapHandler(w http.ResponseWriter, r *http.Request) {
	data, err := LoadBucket("china")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	users := LoadUser(data)
	cm := make(map[string]int)
	for _, city := range CityList {
		cm[city] = 0
	}
	for _, user := range users {
		if user.Location == "" {
			user.Location = "UnKnow"
			continue
		}
		location := user.Location
		city, err := SearchCityFromLocal(location)
		if err == nil {
			cm[city] = cm[city] + 1
			continue
		}
		log.Println("SearchCityFromLocal failed", location, err)
		city, err = SearchCityFromInternet(location)
		if err == nil {
			cm[city] = cm[city] + 1
			continue
		}
		log.Println("SearchCityFromInternet failed", location, err)
	}
	var lts []Location
	for k, v := range cm {
		lts = append(lts, Location{
			Name:   k,
			Amount: v,
		})
	}

	body, err := json.MarshalIndent(map[string]interface{}{
		"status":    0,
		"locations": lts,
	}, "", "    ")
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(body)
}

func WorldMapHandler(w http.ResponseWriter, r *http.Request) {
	data, err := LoadBucket("world")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	users := LoadUser(data)
	cm := make(map[string]int)
	for _, country := range CountryCodeList {
		cm[country] = 0
	}

	for _, user := range users {
		if user.Location == "" {
			user.Location = "UnKnow"
			continue
		}
		location := user.Location
		country, err := SearchCountryFromLocal(location)
		if err == nil {
			cm[country] = cm[country] + 1
			continue
		}
		log.Println("SearchCountryFromLocal failed", location, err)
		country, err = SearchCountryFromInternet(location)
		if err == nil {
			cm[country] = cm[country] + 1
			continue
		}
		log.Println("SearchCountryFromInternet failed", location, err)
	}
	var lts []Location
	for k, v := range cm {
		lts = append(lts, Location{
			Name:   k,
			Amount: v,
		})
	}

	body, err := json.MarshalIndent(map[string]interface{}{
		"status":    0,
		"locations": lts,
	}, "", "    ")
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Write(body)
}
