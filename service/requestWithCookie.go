package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetWithCookieJar(urlString string, jar http.CookieJar) ([]byte, error) {
	client := http.Client{Jar: jar}
	response, err := client.Get(urlString)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(response.Body)
}

func PostJsonWithCookieJar(urlString string, content []byte, jar http.CookieJar) ([]byte, error) {
	client := http.Client{Jar: jar}
	response, err := client.Post(urlString, "application/json", bytes.NewBuffer(content))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(response.Body)
}

func PostFormWithCookieJar(urlString string, content map[string]string, jar http.CookieJar) ([]byte, error) {
	client := http.Client{Jar: jar}
	data := url.Values{}
	for k, v := range content {
		data.Add(k, v)
	}
	fmt.Println(urlString)
	fmt.Println(data)
	response, err := client.PostForm(urlString, data)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(response.Body)
}
