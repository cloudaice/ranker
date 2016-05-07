package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

const (
	Unknow string = "Unknow"
)

func main() {
	logfile, err := os.OpenFile("ranker.log", os.O_WRONLY|os.O_CREATE, 06444)
	if err != nil {
		panic(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	go RunJobs(CrawlChina, storeUser, etlCountry, etlCity, etlScore, etlLanguage)
	go RunJobs(CrawlWorld, storeUser, etlCountry, etlCity, etlScore, etlLanguage)
	RunServer()
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
	log.Println("Start Service on :9090")
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

func LoadUsers(data [][]byte) []User {
	var users []User
	for _, val := range data {
		user := User{}
		if err := json.Unmarshal(val, &user); err != nil {
			log.Printf("Unmarshal user error: %s, user: %s\n", err, string(val))
			continue
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
	users := LoadUsers(data)
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
	users := LoadUsers(data)
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
	users := LoadUsers(data)
	cm := make(map[string]int)
	for _, city := range CityList {
		cm[city] = 0
	}
	for _, user := range users {
		city := user.RealCity
		if city == Unknow {
			continue
		}
		cm[city] = cm[city] + 1
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
	users := LoadUsers(data)
	cm := make(map[string]int)
	for _, country := range CountryCodeList {
		cm[country] = 0
	}

	for _, user := range users {
		country := user.RealCountry
		if country == Unknow {
			continue
		}
		cm[country] = cm[country] + 1
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
