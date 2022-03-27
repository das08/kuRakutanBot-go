package module

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
)

func CreateLINEBotClient(e *Environments) *linebot.Client {
	client, err := linebot.New(
		e.LINE_CHANNEL_SECRET,
		e.LINE_CHANNEL_ACCESS_TOKEN,
	)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
