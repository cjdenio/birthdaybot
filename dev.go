package main

import (
	botapi "birthdaybot/api"
	"birthdaybot/api/slack"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

type handler struct{}

func (h handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if regexp.MustCompile(`(?i)\/api\/slack\/command\/?`).MatchString(req.URL.Path) {
		slack.CommandHandler(res, req)
	} else if regexp.MustCompile(`(?i)\/api\/slack\/events\/?`).MatchString(req.URL.Path) {
		slack.EventsHandler(res, req)
	} else if regexp.MustCompile(`(?i)\/api\/cron\/?`).MatchString(req.URL.Path) {
		botapi.CronHandler(res, req)
	} else {
		res.WriteHeader(404)
		res.Write([]byte("<h2>404: page not found</h2>\nBirthdayBot's dev script doesn't serve static assets, as well as TypeScript code. Please use <code>vercel dev</code> for a more complete experience."))
	}
}

func main() {
	server := &http.Server{
		Addr:    "0.0.0.0:3000",
		Handler: handler{},
	}

	fmt.Println("Dev server was started on port 3000. (unless you see an error below)")

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Something went wrong whilst starting the dev server. :(\nPlease see the error above.")
		os.Exit(1)
	}
}
