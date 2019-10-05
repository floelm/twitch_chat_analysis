package domain

import "twitch_chat_analysis/pkg"

type MessageProcessor struct {
	queuer pkg.Queuer
}

func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{
		queuer: pkg.NewQueuer(),
	}
}

func (p *MessageProcessor) Process() {
	myHandler := func(msg string) {
		println(msg)
	}

	p.queuer.Receive(myHandler)
}
