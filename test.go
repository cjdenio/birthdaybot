package main

import (
	"net/http"
	"birthdaybot/api/slack"
)

type handler struct {}

func (h handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	slack.Handler(res, req)
}

func main() {
	server := http.Server{
		Addr: ":3000",
		Handler: handler{},
	}
	server.ListenAndServe()
}