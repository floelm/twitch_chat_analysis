package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/gempir/go-twitch-irc"
	"github.com/streadway/amqp"
	"log"
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

type Queuer struct {
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewQueuer() Queuer {
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

	return Queuer{
		channel: ch,
		queue:   &q,
	}
}

func (q *Queuer) Produce(message twitch.PrivateMessage) {
	fmt.Println(message.User.Name + " " + message.Message)

	msg := q.buildMessage(message)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	err = q.channel.Publish(
		"",
		q.queue.Name,
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

func (q *Queuer) Receive(handler func(msg string)) {
	msgs, err := q.channel.Consume(
		q.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			handler(string(d.Body))
		}
	}()

	log.Printf("Waiting for messages...")
	<-forever
}

func (q *Queuer) buildMessage(message twitch.PrivateMessage) Message {
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
