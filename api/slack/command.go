package slack

import (
	lib "birthdaybot/api/_lib"
	"birthdaybot/api/_lib/commands"
	"io/ioutil"
	"os"

	"net/http"
	"net/url"
	"strings"
)

// CommandHandler handles the /birthday slash command
func CommandHandler(res http.ResponseWriter, req *http.Request) {
	// Panic handler
	defer func() {
		if err := recover(); err != nil {
			res.Write([]byte("An unknown error occurred."))
		}
	}()

	body, _ := ioutil.ReadAll(req.Body)

	if !lib.SlackVerify(body, os.Getenv("SLACK_SIGNING_SECRET"), req.Header.Get("x-slack-request-timestamp"), req.Header.Get("x-slack-signature")) {
		res.WriteHeader(401)
		res.Write([]byte("Couldn't verify Slack request."))
		return
	}

	parsedBody, _ := url.ParseQuery(string(body))
	req.Form = parsedBody

	switch strings.TrimSpace(strings.ToLower(req.Form["text"][0])) {
	case "":
		commands.HandleDefaultCommand(res, req)
	case "help":
		commands.HandleHelpCommand(res, req)
	case "forget":
		commands.HandleForgetCommand(res, req)
	default:
		commands.HandleHelpCommand(res, req)
	}
}
