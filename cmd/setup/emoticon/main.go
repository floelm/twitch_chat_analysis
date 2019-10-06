package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"twitch_chat_analysis/pkg"
)

func main() {
	client := http.Client{}
	cache := pkg.NewEmoticonsCache()

	maxPages := 2

	for i := 1; i <= maxPages; i++ {
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.frankerfacez.com/v1/emoticons?page=%d", i), nil)
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)

		response := Response{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			panic(err)
		}

		maxPages = response.Pages

		for _, emote := range response.Emoticons {
			cache.Increase(emote.Name)
		}

		println(fmt.Sprintf("Page: %d", i))
	}
}

type Response struct {
	Links struct {
		Next string `json:"next"`
		Self string `json:"self"`
	} `json:"_links"`
	Pages     int `json:"_pages"`
	Total     int `json:"_total"`
	Emoticons []struct {
		Name string `json:"name"`
	} `json:"emoticons"`
}
