package domain

import (
	"strings"
	"twitch_chat_analysis/pkg"
)

type MessageProcessor struct {
	queuer        pkg.Queuer
	termsCache    TermsCache
	emoticonCache pkg.EmoticonsCache
}

func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{
		queuer:        pkg.NewQueuer(),
		termsCache:    NewTermsCache(),
		emoticonCache: pkg.NewEmoticonsCache(),
	}
}

func (p *MessageProcessor) Process() {
	myHandler := func(msg pkg.Message) {
		simpleTerms := strings.Split(msg.Text, " ")

		for _, term := range simpleTerms {
			if p.TermIsEmote(term, msg) {
				continue
			}

			p.termsCache.IncreaseTermCount(term)
		}
	}

	p.queuer.Receive(myHandler)
}

func (p *MessageProcessor) TermIsEmote(term string, msg pkg.Message) bool {
	termIsEmote := false

	for _, emoticon := range msg.Emotes {
		if emoticon.Name == term {
			termIsEmote = true
		}
	}

	if p.emoticonCache.Exists(term) {
		termIsEmote = true
	}

	return termIsEmote
}
