package domain

import (
	"github.com/avast/retry-go"
	"github.com/gempir/go-twitch-irc"
	"twitch_chat_analysis/pkg"
)

type Stalker struct {
	credentials Credentials
	queuer      pkg.Queuer
}

type Credentials struct {
	UserName   string
	OAuthToken string
}

func NewStalker(credentials Credentials) Stalker {
	queuer := pkg.NewQueuer()

	return Stalker{
		credentials: credentials,
		queuer:      queuer,
	}
}

func (s *Stalker) Read(channel string) {
	client := twitch.NewClient(s.credentials.UserName, s.credentials.OAuthToken)
	client.Join(channel)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		s.queuer.Produce(message)
	})

	err := retry.Do(client.Connect)
	if err != nil {
		panic(err)
	}
}
