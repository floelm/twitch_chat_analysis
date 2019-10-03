package domain

import (
	"github.com/gempir/go-twitch-irc"
)

type Stalker struct {
	credentials Credentials
	producer    Producer
}

type Credentials struct {
	UserName   string
	OAuthToken string
}

func NewStalker(credentials Credentials) Stalker {
	producer := NewProducer()

	return Stalker{
		credentials: credentials,
		producer:    producer,
	}
}

func (s *Stalker) Read(channel string) {
	client := twitch.NewClient(s.credentials.UserName, s.credentials.OAuthToken)
	client.Join(channel)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		s.producer.Produce(message)
	})

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
