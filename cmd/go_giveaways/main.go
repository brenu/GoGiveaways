package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/brenu/GoGiveaways/internal/games"
	twitterwrapper "github.com/brenu/GoGiveaways/internal/twitter"
	"github.com/joho/godotenv"

	"github.com/dghubble/go-twitter/twitter"
)

func handleNewGiveaways(lastGiveaway games.GiveAway, wrapper *twitterwrapper.TwitterWrapper) games.GiveAway {
	newGiveaways, err := games.GamesLookUp()

	if err != nil {
		log.Printf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		return lastGiveaway
	}

	newLastGiveaway := newGiveaways[0]

	if lastGiveaway.ID == 0 {
		err := postNewGiveaway(wrapper, newLastGiveaway)

		if err != nil {
			return lastGiveaway
		}
	} else if lastGiveaway.ID != newGiveaways[0].ID {
		newGiveawaysLength := len(newGiveaways)

		for i := 0; i < newGiveawaysLength && newGiveaways[i].ID > lastGiveaway.ID; i++ {
			newLastGiveaway := newGiveaways[i]

			err := postNewGiveaway(wrapper, newLastGiveaway)

			if err != nil {
				return lastGiveaway
			}
		}
	}

	return newGiveaways[0]
}

func postNewGiveaway(wrapper *twitterwrapper.TwitterWrapper, newLastGiveaway games.GiveAway) error {
	imageID := wrapper.HandleImagePost(newLastGiveaway.Image)

	tweetString := fmt.Sprintf("%s - %s\n\nAvailable on: %s", newLastGiveaway.Title, newLastGiveaway.Platforms, newLastGiveaway.GamerpowerURL)

	_, _, err := wrapper.TwitterClient.Statuses.Update(tweetString, &twitter.StatusUpdateParams{MediaIds: []int64{imageID}})

	if err != nil {
		log.Printf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load(filepath.Join("./", ".env"))

	if err != nil {
		log.Printf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		os.Exit(-1)
	}

	twitterWrapper := twitterwrapper.NewTwitterWrapper()

	var lastGiveaway games.GiveAway
	lastGiveaway = handleNewGiveaways(lastGiveaway, twitterWrapper)

	for {
		select {
		case <-time.After(time.Minute):
			lastGiveaway = handleNewGiveaways(lastGiveaway, twitterWrapper)
		}
	}
}
