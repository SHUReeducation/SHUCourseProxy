package model

import (
	"SHUCourseProxy/infrastructure"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type serializedJar struct {
	Cookies map[string][]*http.Cookie
}

func GetCookieJar(studentId string, siteId int16) (http.CookieJar, error) {
	row := infrastructure.DB.QueryRow(`
			SELECT cookie
			FROM cookies
			WHERE student_id=$1 AND site_id=$2;
		`, studentId, siteId)
	var resultBytes []byte
	err := row.Scan(&resultBytes)
	if err != nil {
		return nil, err
	}
	jar, _ := cookiejar.New(nil)
	var serialized serializedJar
	err = json.Unmarshal(resultBytes, &serialized)
	for urlString, cookies := range serialized.Cookies {
		urlObject, _ := url.Parse(urlString)
		jar.SetCookies(urlObject, cookies)
	}
	return jar, err
}

func SetCookieJar(studentId string, siteId int16, jar http.CookieJar) {
	oathUrl, _ := url.Parse("https://oauth.shu.edu.cn/oauth/")
	ssoUrl, _ := url.Parse("https://sso.shu.edu.cn/idp/")
	bytes, _ := json.Marshal(serializedJar{
		Cookies: map[string][]*http.Cookie{
			oathUrl.String(): jar.Cookies(oathUrl),
			ssoUrl.String():  jar.Cookies(ssoUrl),
		},
	})
	_, _ = infrastructure.DB.Exec(`
	DELETE FROM cookies
	WHERE (student_id=$1 AND site_id=$2);
	`, studentId, siteId)
	_, _ = infrastructure.DB.Exec(`
	INSERT INTO cookies(student_id, site_id, cookie)
			VALUES ($1,$2,$3);
	`, studentId, siteId, bytes)
}
