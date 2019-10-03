package main

import (
	"github.com/tkanos/gonfig"
	"sync"
	"twitch_chat_analysis/cmd/reader/domain"
)

type Configuration struct {
	UserName       string
	OAuthToken     string
	ChannelsToJoin []string
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

	var wg sync.WaitGroup
	wg.Add(len(configuration.ChannelsToJoin))

	for _, channel := range configuration.ChannelsToJoin {
		go stalker.Read(channel)
	}

	wg.Wait()
}
