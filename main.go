package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func handleNewGiveaways(lastGiveaway GiveAway, twitterClient *twitter.Client) GiveAway {
	newGiveaways := GamesLookUp()

	if lastGiveaway.ID == 0 {
		newLastGiveaway := newGiveaways[len(newGiveaways)-1]

		tweetString := fmt.Sprintf("%s - %s\n\nAvailable on: %s", newLastGiveaway.Title, newLastGiveaway.Platforms, newLastGiveaway.GamerpowerURL)

		_, _, err := twitterClient.Statuses.Update(tweetString, &twitter.StatusUpdateParams{})

		if err != nil {
			log.Fatalf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		}
	} else if lastGiveaway.ID != newGiveaways[len(newGiveaways)-1].ID {
		newGiveawaysLength := len(newGiveaways)

		for i := 0; i < newGiveawaysLength && newGiveaways[i].ID != lastGiveaway.ID; i++ {
			newLastGiveaway := newGiveaways[i]

			tweetString := fmt.Sprintf("%s - %s\n\nAvailable on: %s", newLastGiveaway.Title, newLastGiveaway.Platforms, newLastGiveaway.GamerpowerURL)

			_, _, err := twitterClient.Statuses.Update(tweetString, &twitter.StatusUpdateParams{})

			if err != nil {
				log.Fatalf(fmt.Sprint("Bad news here, reason: ", err.Error()))
			}
		}
	}

	return newGiveaways[len(newGiveaways)-1]
}

func main() {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	var lastGiveaway GiveAway
	lastGiveaway = handleNewGiveaways(lastGiveaway, client)

	for {
		select {
		case <-time.After(time.Minute):
			lastGiveaway = handleNewGiveaways(lastGiveaway, client)
		}
	}
}
