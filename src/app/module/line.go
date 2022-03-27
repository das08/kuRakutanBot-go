package module

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
)

type LINEBot struct {
	Bot        *linebot.Client
	replyToken string
}

func CreateLINEBotClient(e *Environments) *LINEBot {
	bot, err := linebot.New(
		e.LINE_CHANNEL_SECRET,
		e.LINE_CHANNEL_ACCESS_TOKEN,
	)
	if err != nil {
		log.Fatal(err)
	}
	lb := LINEBot{Bot: bot}
	return &lb
}

func (lb *LINEBot) SetReplyToken(replyToken string) {
	lb.replyToken = replyToken
}

func (lb *LINEBot) SendTextMessage(text string) {
	_, err := lb.Bot.ReplyMessage(lb.replyToken, linebot.NewTextMessage(text)).Do()
	if err != nil {
		log.Print(err)
	}
}

func (lb *LINEBot) SendFlexMessage(flex []byte, altText string) {
	flexContainer, _ := linebot.UnmarshalFlexMessageJSON(flex)
	_, err := lb.Bot.ReplyMessage(lb.replyToken, linebot.NewFlexMessage(altText, flexContainer)).Do()
	if err != nil {
		log.Print(err)
	}
}
