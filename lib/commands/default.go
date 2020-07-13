package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/slack-go/slack"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// HandleDefaultCommand handles the /birthday command with no parameters.
func HandleDefaultCommand(res http.ResponseWriter, req *http.Request) {
	db, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_URL")))
	if err != nil {
		fmt.Println(err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := db.Connect(ctx); err != nil {
		fmt.Println(err.Error())
	}

	defer db.Disconnect(ctx)

	api := slack.New(os.Getenv("SLACK_TOKEN"))

	user := req.Form["user_id"][0]

	profile, err := api.GetUserProfile(user, false)

	rawBirthday, ok := profile.Fields.ToMap()["XfQN2QL49W"]

	if !ok {
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

	parsedDate, err := time.Parse("2006-01-02", rawBirthday.Value)
	if err != nil {
		fmt.Println(err.Error())
	}
	formattedDate := parsedDate.Format("January 2, 2006")

	collection := db.Database("birthdaybot").Collection("birthdays")

	result, err := collection.UpdateOne(ctx, bson.D{{Key: "user_id", Value: user}}, bson.M{"$set": bson.M{"user_id": user, "birthday": rawBirthday.Value, "date": parsedDate.Format("01-02")}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Fatal(err)
	}

	var m map[string]interface{}

	if result.ModifiedCount > 0 || result.UpsertedCount > 0 {
		m = map[string]interface{}{
			"response_type": "ephemeral",
			"blocks": []interface{}{
				map[string]interface{}{
					"type": "section",
					"text": map[string]string{
						"type": "mrkdwn",
						"text": fmt.Sprintf("Thanks, <@%s>! I've remembered your birthday *(%s)*, and I'll post something in <#C0266FRGV> when it comes around! :tada:", user, formattedDate),
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
	} else {
		m = map[string]interface{}{
			"response_type": "ephemeral",
			"blocks": []interface{}{
				map[string]interface{}{
					"type": "section",
					"text": map[string]string{
						"type": "mrkdwn",
						"text": fmt.Sprintf("I already know what your birthday is; it's %s! :birthday:", formattedDate),
					},
				},
			},
		}
	}

	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	res.Header().Set("Content-type", "application/json")
	res.Write(b)
}
