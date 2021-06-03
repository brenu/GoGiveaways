package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type TwitterImage struct {
	MediaID          int64  `json:"media_id"`
	MediaIDString    string `json:"media_id_string"`
	MediaKey         string `json:"media_key"`
	Size             int    `json:"size"`
	ExpiresAfterSecs int    `json:"expires_after_secs"`
	Image            struct {
		ImageType string `json:"image_type"`
		W         int    `json:"w"`
		H         int    `json:"h"`
	} `json:"image"`
}

func handleNewGiveaways(lastGiveaway GiveAway, httpClient *http.Client) GiveAway {
	twitterClient := twitter.NewClient(httpClient)
	newGiveaways, err := GamesLookUp()

	if err != nil {
		log.Fatalf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		return lastGiveaway
	}

	newLastGiveaway := newGiveaways[0]

	if lastGiveaway.ID == 0 {

		imageID := handleImagePost(newLastGiveaway.Image, httpClient)

		tweetString := fmt.Sprintf("%s - %s\n\nAvailable on: %s", newLastGiveaway.Title, newLastGiveaway.Platforms, newLastGiveaway.GamerpowerURL)

		_, _, err := twitterClient.Statuses.Update(tweetString, &twitter.StatusUpdateParams{MediaIds: []int64{imageID}})

		if err != nil {
			log.Fatalf(fmt.Sprint("Bad news here, reason: ", err.Error()))
			return lastGiveaway
		}
	} else if lastGiveaway.ID != newGiveaways[0].ID {
		newGiveawaysLength := len(newGiveaways)

		imageID := handleImagePost(newLastGiveaway.Image, httpClient)

		for i := 0; i < newGiveawaysLength && newGiveaways[i].ID != lastGiveaway.ID; i++ {
			newLastGiveaway := newGiveaways[i]

			tweetString := fmt.Sprintf("%s - %s\n\nAvailable on: %s", newLastGiveaway.Title, newLastGiveaway.Platforms, newLastGiveaway.GamerpowerURL)

			_, _, err := twitterClient.Statuses.Update(tweetString, &twitter.StatusUpdateParams{MediaIds: []int64{imageID}})

			if err != nil {
				log.Fatalf(fmt.Sprint("Bad news here, reason: ", err.Error()))
				return lastGiveaway
			}
		}
	}

	return newGiveaways[0]
}

func handleImagePost(imageURL string, httpClient *http.Client) int64 {
	res, err := http.Get(imageURL)

	if err != nil || res.StatusCode != 200 {
		log.Fatalf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		return 0
	}
	defer res.Body.Close()
	m, _, err := image.Decode(res.Body)
	if err != nil {
		log.Fatalf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		return 0
	}

	form := url.Values{}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	bytesImage := buf.Bytes()

	encodedImage := base64.StdEncoding.EncodeToString(bytesImage)

	form.Add("media_data", encodedImage)

	resp, err := httpClient.PostForm("https://upload.twitter.com/1.1/media/upload.json?media_category=tweet_image", form)

	if err != nil {
		log.Fatalf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		return 0
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		return 0
	}

	var imageResponse TwitterImage

	err = json.Unmarshal(body, &imageResponse)

	return int64(imageResponse.MediaID)
}

func main() {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	var lastGiveaway GiveAway
	lastGiveaway = handleNewGiveaways(lastGiveaway, httpClient)

	for {
		select {
		case <-time.After(time.Minute):
			lastGiveaway = handleNewGiveaways(lastGiveaway, httpClient)
		}
	}
}
