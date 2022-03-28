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

func (lb *LINEBot) SendFlexMessage(flexMessages []FlexMessage) {
	var messages []linebot.SendingMessage
	for _, fm := range flexMessages {
		messages = append(messages, linebot.NewFlexMessage(fm.AltText, fm.FlexContainer))
	}
	_, err := lb.Bot.ReplyMessage(lb.replyToken, messages...).Do()
	if err != nil {
		log.Print(err)
	}
}
