package slack

import (
	lib "birthdaybot/api/_lib"

	"context"
	"fmt"
	"os"
	"time"

	"encoding/json"
	"io/ioutil"

	"net/http"
	"net/url"

	"github.com/slack-go/slack"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CommandHandler handles the /birthday slash command
func CommandHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println(os.Getenv("DB_URL"))

	db, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_URL")))
	if err != nil {
		fmt.Println(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := db.Connect(ctx); err != nil {
		fmt.Println(err.Error())
	}

	body, _ := ioutil.ReadAll(req.Body)

	if !lib.SlackVerify(body, os.Getenv("SLACK_SIGNING_SECRET"), req.Header.Get("x-slack-request-timestamp"), req.Header.Get("x-slack-signature")) {
		res.WriteHeader(401)
		res.Write([]byte("Couldn't verify Slack request."))
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
	rawBirthday, ok := profile.Fields.ToMap()["XfQN2QL49W"]

	if ok {
		date, err := time.Parse("2006-01-02", rawBirthday.Value)
		if err != nil {
			fmt.Println(err.Error())
		}
		parsedDate = date.Format("January 2, 2006")
	} else {
		m := map[string]interface{}{
			"response_type": "ephemeral",
			"blocks": []interface{}{
				map[string]interface{}{
					"type": "section",
					"text": map[string]string{
						"type": "mrkdwn",
						"text": "Hmm... :thinking_face: I'm not sure when your birthday is. Have you set it in your <https://slack.com/help/articles/204092246-Edit-your-profile|profile settings>?",
					},
				},
			},
		}
		b, _ := json.Marshal(m)

		res.Header().Add("Content-type", "application/json")
		res.Write(b)
		return
	}

	collection := db.Database("birthdaybot").Collection("birthdays")

	collection.UpdateOne(ctx, bson.D{{Key: "user_id", Value: user}}, bson.M{"$set": bson.M{"user_id": user, "birthday": rawBirthday.Value}}, options.Update().SetUpsert(true))

	m := map[string]interface{}{
		"response_type": "ephemeral",
		"blocks": []interface{}{
			map[string]interface{}{
				"type": "section",
				"text": map[string]string{
					"type": "mrkdwn",
					"text": fmt.Sprintf("Thanks, <@%s>! I've remembered your birthday *(%s)*, and I'll post something in <#C0266FRGV> when it comes around! :tada:", user, parsedDate),
				},
			},
			map[string]interface{}{
				"type": "context",
				"elements": []interface{}{
					map[string]string{
						"type": "mrkdwn",
						"text": "Type `/birthday forget` to make me forget your birthday.",
					},
				},
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
