package main

import "twitch_chat_analysis/cmd/consumer/domain"

func main() {
	processor := domain.NewMessageProcessor()
	processor.Process()
}
