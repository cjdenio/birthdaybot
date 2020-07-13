package lib

import (
	"net/url"

	"fmt"
	"sort"
	"time"

	"github.com/slack-go/slack"
	"go.mongodb.org/mongo-driver/bson"
)

// GenerateImageURL generate an image URL.
func GenerateImageURL(name string, image string, date string) string {
	thing := url.URL{
		Scheme: "https",
		Host:   "hackclub-birthday-bot.now.sh",
		Path:   "/api/image",
	}
	q := url.Values{}
	q.Set("text", name)
	q.Set("image", image)
	q.Set("date", date)
	thing.RawQuery = q.Encode()

	marshalled, _ := thing.MarshalBinary()

	return string(marshalled)
}

// BirthdaysToBlocks turns an array of birthdays into an array of Block Kit blocks
func BirthdaysToBlocks(birthdays []bson.M) []slack.Block {
	sort.Slice(birthdays, func(i int, j int) bool {
		timeA, _ := time.Parse("01-02", birthdays[i]["date"].(string))
		timeB, _ := time.Parse("01-02", birthdays[j]["date"].(string))

		return timeA.Before(timeB)
	})

	upcomingBlocks := []slack.Block{
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", ":birthday: *Upcoming Birthdays*", false, false), nil, nil),
	}
	pastBlocks := []slack.Block{
		slack.NewDividerBlock(),
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", ":cake: *Past Birthdays*", false, false), nil, nil),
	}

	for _, v := range birthdays {
		parsed, _ := time.Parse("01-02", v["date"].(string))
		_, nowMonth, nowDay := time.Now().Date()
		formatted := parsed.Format("January 2")

		if (parsed.Month() > nowMonth) || (parsed.Month() == nowMonth && parsed.Day() > nowDay) {
			upcomingBlocks = append(upcomingBlocks, slack.NewSectionBlock(nil, []*slack.TextBlockObject{
				slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<@%s>", v["user_id"]), false, false),
				slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*%s*", formatted), false, false),
			}, nil))
		} else {
			pastBlocks = append(pastBlocks, slack.NewSectionBlock(nil, []*slack.TextBlockObject{
				slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("<@%s>", v["user_id"]), false, false),
				slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*%s*", formatted), false, false),
			}, nil))
		}
	}

	return append(upcomingBlocks, pastBlocks...)
}
