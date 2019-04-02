package main

import (
	"SHUCourseProxy/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello? Am I working?"))
	})
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/get", handler.GetWithCookieHandler)
	http.HandleFunc("/post", handler.PostWithCookieHandler)
	_ = http.ListenAndServe(":8086", nil)
}
