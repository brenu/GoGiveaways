package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	// Send a Tweet
	tweet, resp, err := client.Statuses.Update("just setting up my twttr", nil)

	if err != nil {
		log.Fatalf(fmt.Sprint("Bad news here, reason: ", err.Error()))
	}

	log.Output(1, fmt.Sprintf(tweet.FullText, resp))
}
