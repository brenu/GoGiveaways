package games

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// GiveAway is the struct that represents the API response default pattern
type GiveAway struct {
	ID              int64  `json:"id"`
	Title           string `json:"title"`
	Worth           string `json:"worth"`
	Thumbnail       string `json:"thumbnail"`
	Image           string `json:"image"`
	Description     string `json:"description"`
	Instructions    string `json:"instructions"`
	OpenGiveawayURL string `json:"open_giveaways_url"`
	PublishedDate   string `json:"published_date"`
	Type            string `json:"type"`
	Platforms       string `json:"platforms"`
	EndDate         string `json:"end_date"`
	Users           int64  `json:"users"`
	Status          string `json:"status"`
	GamerpowerURL   string `json:"gamerpower_url"`
	OpenGiveaway    string `json:"open_giveaway"`
}

// GamesLookUp is the method that consumes the GamerPower API looking for new free games
func GamesLookUp() ([]GiveAway, error) {
	resp, err := http.Get("https://www.gamerpower.com/api/giveaways")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var giveawaysResponse []GiveAway

	err = json.Unmarshal(body, &giveawaysResponse)

	if err != nil {
		return nil, err
	}

	return giveawaysResponse, nil
}
