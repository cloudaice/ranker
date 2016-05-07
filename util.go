package main

import (
	"errors"
	"strings"
)

var ErrNotFound = errors.New("Not Matched")

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

// 在本地数组中查找匹配的城市，返回城市拼音
func matchCityFromLocal(search string) (string, error) {
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
	return "", ErrNotFound
}

// 在本地数组中寻找匹配的国家，返回国家代码
func matchCountryFromLocal(country string) (string, error) {
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
	return "", ErrNotFound
}
