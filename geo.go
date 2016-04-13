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
	resp, err := client.Get(searchUrl)
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

// SearchCityFromInternet return the best match city from geo website
func SearchCityFromInternet(city string) (string, error) {
	js, err := FetchGeo(city)
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
			log.Println("Change adminName1 error", js.GetIndex(idx).Get("adminName1"), err)
			break
		}
		for _, cty := range CityList {
			if strings.Contains(strings.ToLower(val), cty) {
				return cty, nil
			}
		}
	}
	return "", fmt.Errorf("None")
}

func SearchCountryFromInternet(country string) (string, error) {
	js, err := FetchGeo(country)
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
			log.Println("Change countryCode error", js.GetIndex(idx).Get("countryCode"), err)
			break
		}
		for _, countryCode := range CountryCodeList {
			if strings.Contains(val, countryCode) {
				return countryCode, nil
			}
		}
	}
	return "", fmt.Errorf("None")
}
