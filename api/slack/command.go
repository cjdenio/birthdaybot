package slack

import (
	//"io/ioutil"
	"fmt"
	"net/http"

	//"net/url"
	lib "birthdaybot/api/_lib"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/slack-go/slack"
)

// Handler for slash command
func Handler(res http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)

	if !lib.SlackVerify(body, os.Getenv("SLACK_SIGNING_SECRET"), req.Header.Get("x-slack-request-timestamp"), req.Header.Get("x-slack-signature")) {
		res.WriteHeader(401)
		res.Write([]byte{})
		fmt.Println("Bad")
		return
	}

	api := slack.New(os.Getenv("SLACK_TOKEN"))
	parsedBody, _ := url.ParseQuery(string(body))
	req.Form = parsedBody

	var user string
	if req.Form["text"][0] != "" {
		user = req.Form["text"][0]
	} else {
		user = req.Form["user_id"][0]
	}

	profile, err := api.GetUserProfile(user, false)

	var parsedDate string

	if d, ok := profile.Fields.ToMap()["XfQN2QL49W"]; ok {
		date, err := time.Parse("2006-01-02", d.Value)
		if err != nil {
			fmt.Println(err.Error())
		}
		parsedDate = date.Format("January 2, 2006")
	} else {
		parsedDate = "an unknown date"
	}

	m := map[string]interface{}{
		"response_type": "in_channel",
		"blocks": []interface{}{
			map[string]interface{}{
				"type": "section",
				"text": map[string]string{
					"type": "mrkdwn",
					"text": fmt.Sprintf("Hello <@%s>...", user),
				},
			},
			map[string]string{
				"type":      "image",
				"image_url": fmt.Sprintf("https://birthday-bot.cjdenio.now.sh/api/image?text=%s&date=%s&image=%s", profile.DisplayName, parsedDate, profile.Image192),
				"alt_text":  "Happy Birthday!",
			},
		},
	}
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	res.Header().Set("Content-type", "application/json")
	res.Write(b)
}
