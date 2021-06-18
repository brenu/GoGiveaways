package twitterwrapper

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

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// TwitterImage is the struct that represents a default twitter response for media uploading
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

// TwitterWrapper needs two different clients because dghubble's lib is not able to send tweets with pictures
type TwitterWrapper struct {
	TwitterClient *twitter.Client
	HTTPClient    *http.Client
}

func NewTwitterWrapper() *TwitterWrapper {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	twitterClient := twitter.NewClient(httpClient)

	return &TwitterWrapper{
		TwitterClient: twitterClient,
		HTTPClient:    httpClient,
	}
}

func (t *TwitterWrapper) HandleImagePost(imageURL string) int64 {
	res, err := http.Get(imageURL)

	if err != nil || res.StatusCode != 200 {
		log.Printf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		return 0
	}
	defer res.Body.Close()
	m, _, err := image.Decode(res.Body)
	if err != nil {
		log.Printf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		return 0
	}

	form := url.Values{}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	bytesImage := buf.Bytes()

	encodedImage := base64.StdEncoding.EncodeToString(bytesImage)

	form.Add("media_data", encodedImage)

	resp, err := t.HTTPClient.PostForm("https://upload.twitter.com/1.1/media/upload.json?media_category=tweet_image", form)

	if err != nil {
		log.Printf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		return 0
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		return 0
	}

	var imageResponse TwitterImage

	err = json.Unmarshal(body, &imageResponse)

	if err != nil {
		log.Printf(fmt.Sprint("Bad news here, reason: ", err.Error()))
		return 0
	}

	return int64(imageResponse.MediaID)
}
