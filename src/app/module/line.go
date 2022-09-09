package module

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
)

type LINEBot struct {
	Bot         *linebot.Client
	replyToken  string
	senderUid   string
	isMockUser  bool
	mockContext *gin.Context
}

func CreateLINEBotClient(e *Environments, c *gin.Context) *LINEBot {
	bot, err := linebot.New(
		e.LineChannelSecret,
		e.LineChannelAccessToken,
	)
	if err != nil {
		log.Fatal(err)
	}
	lb := LINEBot{Bot: bot, mockContext: c}
	return &lb
}

func (lb *LINEBot) SetReplyToken(replyToken string) {
	lb.replyToken = replyToken
}

func (lb *LINEBot) SetSenderUid(e *Environments, senderUid string) {
	lb.senderUid = senderUid
	if senderUid == e.LineMockUid {
		lb.isMockUser = true
	}
}

func (lb *LINEBot) SendTextMessage2(text string) {
	if lb.isMockUser {
		lb.mockContext.JSON(200, text)
		return
	}
	_, err := lb.Bot.ReplyMessage(lb.replyToken, linebot.NewTextMessage(text)).Do()
	if err != nil {
		log.Print(err)
	}
}

func (lb *LINEBot) SendFlexMessage(flexMessages FlexMessages) {
	if lb.isMockUser {
		lb.mockContext.JSON(200, flexMessages)
		return
	}
	var messages []linebot.SendingMessage
	for _, fm := range flexMessages {
		messages = append(messages, linebot.NewFlexMessage(fm.AltText, fm.FlexContainer))
	}
	_, err := lb.Bot.ReplyMessage(lb.replyToken, messages...).Do()
	if err != nil {
		log.Print(err)
	}
}
