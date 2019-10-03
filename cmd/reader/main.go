package main

import (
	"github.com/tkanos/gonfig"
	"sync"
	"twitch_chat_analysis/cmd/reader/domain"
)

type Configuration struct {
	UserName      string
	OAuthToken    string
	ChannelToJoin string
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
	wg.Add(1)

	go stalker.Read(configuration.ChannelToJoin)

	wg.Wait()
}
