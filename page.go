package main

type PageAttr struct {
	Name      string
	Location  string
	Followers int
}

func PageParser(url string) PageAttr {
	return PageAttr{}
}
