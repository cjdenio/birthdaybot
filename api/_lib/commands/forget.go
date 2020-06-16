package commands

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// HandleForgetCommand handles /birthday forget
func HandleForgetCommand(res http.ResponseWriter, req *http.Request) {
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

	collection := db.Database("birthdaybot").Collection("birthdays")

	result, err := collection.DeleteOne(ctx, bson.D{{Key: "user_id", Value: req.Form["user_id"][0]}})

	if err != nil {
		fmt.Println(err.Error())
	}

	if result.DeletedCount == 0 {
		res.Write([]byte("There's nothing to forget; I don't know what your birthday is! :shrug:"))
	} else {
		res.Write([]byte("Much like a relative, I've successfully forgotten your birthday. :+1:"))
	}
}
