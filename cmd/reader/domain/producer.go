package domain

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc"
	"github.com/streadway/amqp"
	"time"
)

type Message struct {
	Id      string
	User    User
	Text    string
	Channel string
	RoomID  string
	Time    time.Time
	Emotes  []Emote
	Bits    int
	Action  bool
}

type User struct {
	Id          string
	Name        string
	DisplayName string
	Color       string
	Badges      map[string]int
}

type Emote struct {
	Name  string
	Id    string
	Count int
}

type Producer struct {
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewProducer() Producer {
	conn, err := amqp.Dial("amqp://rabbitmquser:some_password@localhost:7001/")
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare(
		"messages",
		false,
		false,
		false,
		false,
		nil,
	)

	return Producer{
		channel: ch,
		queue:   &q,
	}
}

func (p *Producer) Produce(message twitch.PrivateMessage) {
	fmt.Println(message.User.Name + " " + message.Message)

	msg := p.buildMessage(message)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	err = p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgBytes,
		},
	)
	if err != nil {
		panic(err)
	}
}

func (p *Producer) buildMessage(message twitch.PrivateMessage) Message {
	emotes := make([]Emote, 0)

	for _, emote := range message.Emotes {
		emotes = append(emotes, Emote{
			Id:    emote.ID,
			Name:  emote.Name,
			Count: emote.Count,
		})
	}

	return Message{
		Id:   message.ID,
		Text: message.Message,
		User: User{
			Name:        message.User.Name,
			Badges:      message.User.Badges,
			Color:       message.User.Color,
			DisplayName: message.User.DisplayName,
			Id:          message.User.ID,
		},
		Channel: message.Channel,
		RoomID:  message.RoomID,
		Action:  message.Action,
		Time:    message.Time,
		Bits:    message.Bits,
		Emotes:  emotes,
	}
}
