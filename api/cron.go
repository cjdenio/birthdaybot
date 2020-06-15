package api

import (
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"fmt"
	"time"
)

// CronHandler handles the daily cron job.
func CronHandler(res http.ResponseWriter, req *http.Request) {
	//api := slack.New(os.Getenv("SLACK_TOKEN"))

	res.Write([]byte("OK"))

	db, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("DB_URL")))

	if err != nil {
		fmt.Print(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	if err := db.Connect(ctx); err != nil {
		fmt.Println(err.Error())
	}

	collection := db.Database("birthdaybot").Collection("birthdays")
	count, err := collection.CountDocuments(ctx, bson.D{{Key: "birthday", Value: time.Now().Format("2006-01-02")}})

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(count)
}
