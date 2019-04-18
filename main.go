package main

import (
	"SHUCourseProxy/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/get", handler.GetWithCookieHandler)
	http.HandleFunc("/post", handler.PostWithCookieHandler)
	http.HandleFunc("/post-form", handler.PostFormWithCookieHandler)
	_ = http.ListenAndServe(":8086", nil)
}
