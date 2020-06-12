package main

import (
	"birthdaybot/api/slack"
	"net/http"
)

type handler struct{}

func (h handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	slack.CommandHandler(res, req)
}

func main() {
	server := http.Server{
		Addr:    ":3000",
		Handler: handler{},
	}
	server.ListenAndServe()
}
