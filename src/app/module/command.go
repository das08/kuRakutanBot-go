package module

import (
	"github.com/google/uuid"
	"strings"
)

type Command struct {
	Keyword      string
	SendFunction func(c Clients, env *Environments, lb *LINEBot)
}

var Commands = [...]Command{
	{Keyword: "help", SendFunction: helpCmd},
	{Keyword: "へるぷ", SendFunction: helpCmd},
	{Keyword: "ヘルプ", SendFunction: helpCmd},
	{Keyword: "はんてい", SendFunction: judgeDetailCmd},
	{Keyword: "判定", SendFunction: judgeDetailCmd},
	{Keyword: "詳細", SendFunction: judgeDetailCmd},
	{Keyword: "判定詳細", SendFunction: judgeDetailCmd},
	{Keyword: "楽単詳細", SendFunction: judgeDetailCmd},
	{Keyword: "楽単", SendFunction: rakutanCmd},
	{Keyword: "おみくじ", SendFunction: rakutanCmd},
	{Keyword: "楽単おみくじ", SendFunction: rakutanCmd},
	{Keyword: "鬼単", SendFunction: onitanCmd},
	{Keyword: "鬼単おみくじ", SendFunction: onitanCmd},
	{Keyword: "お気に入り", SendFunction: getFavoritesCmd},
	{Keyword: "おきにいり", SendFunction: getFavoritesCmd},
	{Keyword: "リスト", SendFunction: getFavoritesCmd},
	{Keyword: "一覧", SendFunction: getFavoritesCmd},
	{Keyword: "＠info", SendFunction: infoCmd},
	{Keyword: "@info", SendFunction: infoCmd},
	{Keyword: "認証", SendFunction: verificationCmd},
	{Keyword: "認証する", SendFunction: verificationCmd},
	{Keyword: "ユーザ認証", SendFunction: verificationCmd},
	{Keyword: "お問い合わせ", SendFunction: inquiryCmd},
	{Keyword: "問い合わせ", SendFunction: inquiryCmd},
	{Keyword: "myuid", SendFunction: myUIDCmd},
	{Keyword: "京大楽単bot", SendFunction: iconCmd},
}

func IsCommand(messageText string) (bool, func(c Clients, env *Environments, lb *LINEBot)) {
	isCommand := false
	var function func(c Clients, env *Environments, lb *LINEBot)
	for _, cmd := range Commands {
		// Case-insensitive
		if strings.EqualFold(cmd.Keyword, messageText) {
			isCommand = true
			function = cmd.SendFunction
		}
	}
	return isCommand, function
}

func helpCmd(_ Clients, _ *Environments, lb *LINEBot) {
	flexMessages := loadFlexMessages("./assets/richmenu/help.json", "ヘルプ")
	lb.SendFlexMessage(flexMessages)
}

func judgeDetailCmd(_ Clients, _ *Environments, lb *LINEBot) {
	flexMessages := loadFlexMessages("./assets/richmenu/judge_detail.json", "らくたん判定の詳細")
	lb.SendFlexMessage(flexMessages)
}

func inquiryCmd(_ Clients, _ *Environments, lb *LINEBot) {
	flexMessages := loadFlexMessages("./assets/richmenu/inquiry.json", "お問い合わせ")
	lb.SendFlexMessage(flexMessages)
}

func iconCmd(_ Clients, _ *Environments, lb *LINEBot) {
	flexMessages := loadFlexMessages("./assets/richmenu/icon.json", "京大楽単bot")
	lb.SendFlexMessage(flexMessages)
}

func infoCmd(_ Clients, _ *Environments, lb *LINEBot) {
	flexMessages := loadFlexMessages("./assets/richmenu/info.json", "お知らせ")
	lb.SendFlexMessage(flexMessages)
}

func loadFlexMessages(filename string, altText string) []FlexMessage {
	json := LoadJSON(filename)
	return CreateFlexMessage(json, altText)
}

func rakutanCmd(c Clients, env *Environments, lb *LINEBot) {
	status, ok := GetRakutanInfo(c, Omikuji, "rakutan")
	//countUp(env, c.Mongo, lb.senderUid, "rakutan") // TODO
	if ok {
		flexMessages := CreateRakutanDetail(status.Result[0], env, Rakutan)
		lb.SendFlexMessage(flexMessages)
	} else {
		lb.SendTextMessage2(status.Err)
	}
}

func onitanCmd(c Clients, env *Environments, lb *LINEBot) {
	status, ok := GetRakutanInfo(c, Omikuji, "onitan")
	//countUp(env, c.Mongo, lb.senderUid, "onitan") // TODO
	if ok {
		flexMessages := CreateRakutanDetail(status.Result[0], env, Rakutan)
		lb.SendFlexMessage(flexMessages)
	} else {
		lb.SendTextMessage2(status.Err)
	}
}

func getFavoritesCmd(c Clients, env *Environments, lb *LINEBot) {
	queryStatus, ok := c.Postgres.GetFavorites(lb.senderUid)
	if !ok {
		lb.SendTextMessage2(queryStatus.Err)
		return
	}
	flexMessages := CreateFavorites2(queryStatus.Result)
	lb.SendFlexMessage(flexMessages)
}

func verificationCmd(c Clients, env *Environments, lb *LINEBot) {
	verified, err := c.Postgres.IsVerified(lb.senderUid)
	if err != nil {
		lb.SendTextMessage2(ErrorMessageCheckVerificateError)
		return
	}
	var flexMessages []FlexMessage
	if verified {
		flexMessages = loadFlexMessages("./assets/richmenu/verified.json", "ユーザー認証済み")
	} else {
		flexMessages = loadFlexMessages("./assets/richmenu/verification.json", "ユーザー認証をする")
	}
	lb.SendFlexMessage(flexMessages)
}

func myUIDCmd(_ Clients, _ *Environments, lb *LINEBot) {
	lb.SendTextMessage2(lb.senderUid)
}

func SendVerificationCmd(c Clients, env *Environments, lb *LINEBot, email string) {
	uuidObj, _ := uuid.NewUUID()
	data := []byte(lb.senderUid)
	code := uuid.NewSHA1(uuidObj, data).String()
	err := c.Postgres.InsertVerificationToken(lb.senderUid, code)
	if err != nil {
		lb.SendTextMessage2(ErrorMessageInsertVerificateError)
		return
	}
	err = SendVerification(env, email, code)
	if err != nil {
		lb.SendTextMessage2(ErrorMessageVerificationTokenSendError)
		return
	}
	lb.SendTextMessage2(SuccessVericationTokenSent)
}
