package main

func RunJobs(crl Crawler, pst Persister, etls ...Etler) {
	val := func(users <-chan User) <-chan User {
		for _, etl := range etls {
			users = etl(users)
		}
		return users
	}
	pst(val(crl()))
}
