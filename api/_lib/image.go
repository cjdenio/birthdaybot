package lib

import (
	"net/url"
)

// GenerateURL generate an image URL.
func GenerateURL(name string, image string, date string) string {
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
