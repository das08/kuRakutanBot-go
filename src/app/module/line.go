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

type ReplyText struct {
	Status KRBStatus
	Text   string
}

type ReplyFlex struct {
	Status KRBStatus
	Flex   []FlexMessage
}

func CreateLINEBotClient(e *Environments, c *gin.Context) *LINEBot {
	bot, err := linebot.New(
		e.LINE_CHANNEL_SECRET,
		e.LINE_CHANNEL_ACCESS_TOKEN,
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
	if senderUid == e.LINE_MOCK_UID {
		lb.isMockUser = true
	}
}

func (lb *LINEBot) SendTextMessage(rt ReplyText) {
	if lb.isMockUser {
		lb.mockContext.JSON(200, rt)
		return
	}
	_, err := lb.Bot.ReplyMessage(lb.replyToken, linebot.NewTextMessage(rt.Text)).Do()
	if err != nil {
		log.Print(err)
	}
}

func (lb *LINEBot) SendFlexMessage(flexMessages []FlexMessage) {
	if lb.isMockUser {
		lb.mockContext.JSON(200, ReplyFlex{Status: KRBSuccess, Flex: flexMessages})
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
