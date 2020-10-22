package api

import (
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"fmt"
	"time"

	"sync"

	"birthdaybot/lib"

	"github.com/slack-go/slack"
)

// CronHandler handles the daily cron job.
func CronHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Query().Get("token") != os.Getenv("CRON_TOKEN") {
		res.WriteHeader(400)
		res.Write([]byte("Cron request couldn't be verified. :( Nice try, though!"))
		return
	}

	api := slack.New(os.Getenv("SLACK_TOKEN"), slack.OptionDebug(true))

	db, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_URL")))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.Connect(ctx); err != nil {
		log.Fatal(err.Error())
	}
	defer db.Disconnect(ctx)

	collection := db.Database("birthdaybot").Collection("birthdays")

	timezone, _ := time.LoadLocation("America/New_York")
	cursor, err := collection.Find(ctx, bson.D{{Key: "date", Value: time.Now().In(timezone).Format("01-02")}})

	if err != nil {
		log.Fatal(err.Error())
	}

	var results []bson.M
	cursor.All(ctx, &results)

	wg := sync.WaitGroup{}

	wg.Add(len(results))

	for _, v := range results {
		go func(user bson.M) {
			channel := "G014FJELTHP"
			if os.Getenv("GO_ENV") != "development" {
				channel = "C0266FRGV"
			}

			userInfo, err := api.GetUserProfile(user["user_id"].(string), false)

			if err != nil {
				log.Fatal(err)
			}

			userName := userInfo.DisplayName
			if userName == "" {
				userName = userInfo.RealName
			}

			parsedDate, _ := time.Parse("2006-01-02", user["birthday"].(string))
			formattedDate := parsedDate.Format("January 2, 2006")

			if parsedDate.Year() >= 2020 {
				// They probably left out the year (or this code is being run years in the future)
				formattedDate = parsedDate.Format("January 2")
			}

			_, _, err = api.PostMessage(channel, slack.MsgOptionBlocks(
				slack.NewSectionBlock(
					slack.NewTextBlockObject(
						"mrkdwn",
						fmt.Sprintf("It's <@%s>'s birthday! :tada: From all of your fellow Hack Clubbers, have a great one! :partyparrot:", user["user_id"]), false, false), nil, nil),
				slack.NewImageBlock(lib.GenerateImageURL(userName, userInfo.Image192, formattedDate), "Happy birthday!", "image", slack.NewTextBlockObject("plain_text", "Happy birthday!", false, false)),
				slack.NewContextBlock("context", slack.NewTextBlockObject("mrkdwn", "Want me to post something when _your_ special day comes around? Just type `/birthday` to get started!", false, false)),
			),
				slack.MsgOptionText(fmt.Sprintf("It's %s's birthday! From all of your fellow Hack Clubbers, have a great one!", userName), false),
			)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(v)
	}

	wg.Wait()

	res.Write([]byte("cool"))
}
