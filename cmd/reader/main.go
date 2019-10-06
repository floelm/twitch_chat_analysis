package main

import (
	"github.com/tkanos/gonfig"
	"sync"
	"twitch_chat_analysis/cmd/reader/domain"
)

type Configuration struct {
	UserName            string
	OAuthToken          string
	ApplicationClientId string
}

func main() {
	configuration := Configuration{}
	err := gonfig.GetConf("config.dev.json", &configuration)
	if err != nil {
		panic(err)
	}

	stalker := domain.NewStalker(domain.Credentials{
		UserName:   configuration.UserName,
		OAuthToken: configuration.OAuthToken,
	})

	client := domain.NewTwitchClient(configuration.ApplicationClientId)
	channels := client.GetTopChannels()

	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, channel := range channels {
		go stalker.Read(channel)
	}

	wg.Wait()
}
