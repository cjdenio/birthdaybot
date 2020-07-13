package slack

import (
	lib "birthdaybot/lib"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"context"
	"time"

	"github.com/slack-go/slack"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type eventPayload struct {
	Type      string                 `json:"type"`
	Challenge string                 `json:"challenge"`
	Event     map[string]interface{} `json:"event"`
}

// EventsHandler handles events.
func EventsHandler(res http.ResponseWriter, req *http.Request) {
	db, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_URL")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = db.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect(ctx)

	client := slack.New(os.Getenv("SLACK_TOKEN"))

	var body eventPayload
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		log.Fatal(err)
	}
	switch body.Type {
	case "url_verification":
		res.Write([]byte(body.Challenge))
	case "event_callback":
		switch body.Event["type"] {
		case "app_home_opened":
			collection := db.Database("birthdaybot").Collection("birthdays")
			cursor, err := collection.Find(ctx, bson.D{})

			var results []bson.M
			cursor.All(ctx, &results)

			blocks := lib.BirthdaysToBlocks(results)

			_, err = client.PublishView(body.Event["user"].(string), slack.HomeTabViewRequest{
				Type: slack.VTHomeTab,
				Blocks: slack.Blocks{
					BlockSet: blocks,
				},
			}, "")
			if err != nil {
				log.Fatal(err)
			}
		}
		res.Write([]byte("OK"))
	default:
		res.Write([]byte("OK"))
	}
}
