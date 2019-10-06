package domain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type TwitchClient struct {
	applicationClientId string
	client              *http.Client
}

func NewTwitchClient(applicationClientId string) TwitchClient {
	return TwitchClient{
		applicationClientId: applicationClientId,
		client:              &http.Client{},
	}
}

func (c *TwitchClient) GetTopChannels(channelCount int) []string {
	topChannels := make([]string, 0)

	for i := 0; i <= channelCount-100; i = i + 100 {
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/kraken/streams/?limit=100&language=en&offset=%d", i), nil)
		req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
		req.Header.Set("Client-ID", c.applicationClientId)
		resp, err := c.client.Do(req)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)

		channelsResponse := TopChannelsResponse{}
		err = json.Unmarshal(body, &channelsResponse)
		if err != nil {
			panic(err)
		}

		for _, stream := range channelsResponse.Streams {
			topChannels = append(topChannels, stream.Channel.Name)
		}
	}

	return topChannels
}

type TopChannelsResponse struct {
	Total   int `json:"_total"`
	Streams []struct {
		ID      int64 `json:"_id"`
		Channel struct {
			ID                           int         `json:"_id"`
			BroadcasterLanguage          string      `json:"broadcaster_language"`
			CreatedAt                    time.Time   `json:"created_at"`
			DisplayName                  string      `json:"display_name"`
			Followers                    int         `json:"followers"`
			Game                         string      `json:"game"`
			Language                     string      `json:"language"`
			Logo                         string      `json:"logo"`
			Mature                       bool        `json:"mature"`
			Name                         string      `json:"name"`
			Partner                      bool        `json:"partner"`
			ProfileBanner                string      `json:"profile_banner"`
			ProfileBannerBackgroundColor interface{} `json:"profile_banner_background_color"`
			Status                       string      `json:"status"`
			UpdatedAt                    time.Time   `json:"updated_at"`
			URL                          string      `json:"url"`
			VideoBanner                  string      `json:"video_banner"`
			Views                        int         `json:"views"`
		} `json:"channel"`
		CreatedAt  time.Time `json:"created_at"`
		Delay      int       `json:"delay"`
		Game       string    `json:"game"`
		IsPlaylist bool      `json:"is_playlist"`
		Preview    struct {
			Large    string `json:"large"`
			Medium   string `json:"medium"`
			Small    string `json:"small"`
			Template string `json:"template"`
		} `json:"preview"`
		VideoHeight int `json:"video_height"`
		Viewers     int `json:"viewers"`
	} `json:"streams"`
}
