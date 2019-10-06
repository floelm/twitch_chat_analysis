package domain

import (
	"strings"
	"twitch_chat_analysis/pkg"
)

var filterTerms = []string{
	"ㅋ",
	"the",
	"and",
	"!",
}

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

		if p.MessageIsCommand(msg) {
			return
		}

		for _, term := range simpleTerms {
			if p.TermIsEmote(term, msg) {
				continue
			}

			if p.TermIsFiltered(term, msg) {
				continue
			}

			if p.TermIsTooShort(term, msg) {
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

func (p *MessageProcessor) TermIsFiltered(term string, msg pkg.Message) bool {
	termIsFiltered := false

	for _, filterTerm := range filterTerms {
		if strings.Contains(term, filterTerm) {
			termIsFiltered = true
		}
	}

	return termIsFiltered
}

func (p *MessageProcessor) TermIsTooShort(term string, msg pkg.Message) bool {
	// no idea whether this makes sense
	if len(term) <= 3 {
		return true
	}

	return false
}

func (p *MessageProcessor) MessageIsCommand(msg pkg.Message) bool {
	return msg.Action
}
