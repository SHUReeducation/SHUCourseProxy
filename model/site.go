package model

import (
	"SHUCourseProxy/infrastructure"
	"github.com/hashicorp/golang-lru"
	"net/url"
)

var sitesForURLCache *lru.ARCCache

func init() {
	sitesForURLCache, _ = lru.NewARC(512)
}

func GetSiteIdForURL(url string) (int16, error) {
	i, contains := sitesForURLCache.Get(url)
	var id int16
	if !contains {
		row := infrastructure.DB.QueryRow(`
			SELECT id
			FROM site
			WHERE position(domain_name in $1::text) != 0
		`, url)
		err := row.Scan(&id)
		if err != nil {
			return -1, err
		}
		sitesForURLCache.Add(url, id)
	} else {
		id = i.(int16)
	}
	return id, nil
}

func GetOrCreateSiteIdForURL(urlString string) (int16, error) {
	urlObject, _ := url.Parse(urlString)
	alreadyIn, err := GetSiteIdForURL(urlObject.Host)
	if err != nil || alreadyIn == 0 {
		row := infrastructure.DB.QueryRow(`
			INSERT INTO site(domain_name) VALUES ($1) RETURNING id;
		`, urlObject.Host)
		err = row.Scan(&alreadyIn)
	}
	return alreadyIn, err
}
