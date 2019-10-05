package domain

import (
	"strings"
	"twitch_chat_analysis/pkg"
)

type MessageProcessor struct {
	queuer     pkg.Queuer
	termsCache TermsCache
}

func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{
		queuer:     pkg.NewQueuer(),
		termsCache: NewTermsCache(),
	}
}

func (p *MessageProcessor) Process() {
	myHandler := func(msg pkg.Message) {
		simpleTerms := strings.Split(msg.Text, " ")

		for _, term := range simpleTerms {
			p.termsCache.IncreaseTermCount(term)
		}
	}

	p.queuer.Receive(myHandler)
}
