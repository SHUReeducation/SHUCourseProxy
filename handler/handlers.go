package handler

import (
	"SHUCourseProxy/model"
	"SHUCourseProxy/service"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

func getCookieJarFromRequest(r *http.Request, url string) (http.CookieJar, error) {
	tokenString := r.Header.Get("Authorization")[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)
	studentId := claims["studentId"].(string)
	siteId, err := model.GetSiteIdForURL(url)
	if err != nil {
		return nil, err
	}
	jar, err := model.GetCookieJar(studentId, siteId)
	if err != nil {
		return nil, err
	}
	return jar, nil
}

func postWithSaml(urlString string, samlRequest string, relayState string, client *http.Client) (*http.Response, error) {
	return client.PostForm(urlString, url.Values{
		"SAMLRequest": []string{samlRequest},
		"RelayState":  []string{relayState},
	})
}

func simulateLogin(fromURL string, studentId string, password string) http.CookieJar {
	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}
	fmt.Println(fromURL)
	page1, _ := client.Get(fromURL)
	body, _ := ioutil.ReadAll(page1.Body)
	fmt.Println(string(body))
	doc, _ := goquery.NewDocumentFromReader(page1.Body)
	fmt.Println(doc.Html())
	saml, _ := doc.Find("input[name=SAMLRequest]").Attr("value")
	relay, _ := doc.Find("input[name=RelayState]").Attr("value")
	_, _ = postWithSaml("https://sso.shu.edu.cn/idp/profile/SAML2/POST/SSO", saml, relay, &client)
	page2, _ := client.PostForm("https://sso.shu.edu.cn/idp/Authn/UserPassword", url.Values{
		"j_username": []string{studentId},
		"j_password": []string{password},
	})
	doc, _ = goquery.NewDocumentFromReader(page2.Body)
	saml, _ = doc.Find("input[name=SAMLResponse]").Attr("value")
	relay, _ = doc.Find("input[name=RelayState]").Attr("value")
	_, _ = postWithSaml("http://oauth.shu.edu.cn/oauth/Shibboleth.sso/SAML2/POST", saml, relay, &client)
	return client.Jar
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadAll(r.Body)
	var input struct {
		FromUrl  string `json:"from_url"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.Unmarshal(content, &input)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	jar := simulateLogin(input.FromUrl, input.Username, input.Password)
	siteId, _ := model.GetOrCreateSiteIdForURL(input.FromUrl)
	model.SetCookieJar(input.Username, siteId, jar)
	_, _ = w.Write([]byte(service.GenerateJWT(input.Username)))
}

func GetWithCookieHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	var input struct {
		Url string `json:"url"`
	}
	err = json.Unmarshal(body, &input)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	jar, err := getCookieJarFromRequest(r, input.Url)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	result, _ := service.GetWithCookieJar(input.Url, jar)
	_, err = w.Write(result)
	if err != nil {
		w.WriteHeader(500)
	}
}

func PostWithCookieHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	type Content struct {
		Url     string      `json:"url"`
		Content interface{} `json:"content"`
	}
	var content Content
	err = json.Unmarshal(body, &content)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	jar, err := getCookieJarFromRequest(r, content.Url)
	if err != nil {
		w.WriteHeader(403)
		return
	}
	encoded, err := json.Marshal(content.Content)
	result, _ := service.PostJsonWithCookieJar(content.Url, encoded, jar)
	_, err = w.Write(result)
	if err != nil {
		w.WriteHeader(500)
	}
}
