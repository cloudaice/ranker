package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	sjson "github.com/bitly/go-simplejson"
)

func FetchGeo(query string) (*sjson.Json, error) {
	query = url.QueryEscape(query)
	searchUrl := fmt.Sprintf("http://api.geonames.org/searchJSON?q=%s&maxRows=10&username=%s", query, "cloudaice")
	resp, err := httpClient.Get(searchUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	js, err := sjson.NewFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return js, nil
}

// 从网络上匹配对应的城市
func matchCityFromInternet(query string) (string, error) {
	js, err := FetchGeo(query)
	if err != nil {
		return "", err
	}
	numResult, err := js.Get("totalResultsCount").Int()
	if err != nil {
		return "", err
	}

	js = js.Get("geonames")
	for idx := 0; idx < numResult; idx++ {
		val, err := js.GetIndex(idx).Get("adminName1").String()
		if err != nil {
			log.Printf("checkout adminName1 error: %s, adminName1: %v\n", err, js.GetIndex(idx).Get("adminName1"))
			break
		}
		for _, cty := range CityList {
			if strings.Contains(strings.ToLower(val), cty) {
				return cty, nil
			}
		}
	}
	return "", ErrNotFound
}

// 从网络上匹配对应的国家
func matchCountryFromInternet(query string) (string, error) {
	js, err := FetchGeo(query)
	if err != nil {
		return "", err
	}
	numResult, err := js.Get("totalResultsCount").Int()
	if err != nil {
		return "", err
	}
	js = js.Get("geonames")
	for idx := 0; idx < numResult; idx++ {
		val, err := js.GetIndex(idx).Get("countryCode").String()
		if err != nil {
			log.Println("checkout countryCode error: %s, countryCode: %v\n", err, js.GetIndex(idx).Get("countryCode"))
			break
		}
		for _, countryCode := range CountryCodeList {
			if strings.Contains(val, countryCode) {
				return countryCode, nil
			}
		}
	}
	return "", ErrNotFound
}
