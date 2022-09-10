package module

import (
	"github.com/das08/kuRakutanBot-go/richmenu"
	"github.com/google/uuid"
	"log"
	"strings"
)

var (
	JudgeJson        FlexMessages
	InquiryJson      FlexMessages
	IconJson         FlexMessages
	InfoJson         FlexMessages
	VerifiedJson     FlexMessages
	VerificationJson FlexMessages
	HelpJson         FlexMessages
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

func PreloadJson() {
	HelpJson = loadFlexMessages("./assets/richmenu/help.json", "ヘルプ")
	JudgeJson = loadFlexMessages("./assets/richmenu/judge_detail.json", "らくたん判定の詳細")
	InquiryJson = loadFlexMessages("./assets/richmenu/inquiry.json", "お問い合わせ")
	IconJson = loadFlexMessages("./assets/richmenu/icon.json", "京大楽単bot")
	InfoJson = loadFlexMessages("./assets/richmenu/info.json", "お知らせ")
	VerifiedJson = loadFlexMessages("./assets/richmenu/verified.json", "認証済み")
	VerificationJson = loadFlexMessages("./assets/richmenu/verification.json", "認証")
	log.Println("PreloadJson: Done")
	log.Printf("PreloadJson: %v", HelpJson)
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

func helpCmd(c Clients, _ *Environments, lb *LINEBot) {
	go c.Postgres.InsertUserAction(lb.senderUid, UserActionHelp)
	log.Printf("helpCmd: %s %v", lb.senderUid, HelpJson)
	lb.SendFlexMessage(HelpJson)
}

func judgeDetailCmd(_ Clients, _ *Environments, lb *LINEBot) {
	lb.SendFlexMessage(JudgeJson)
}

func inquiryCmd(_ Clients, _ *Environments, lb *LINEBot) {
	lb.SendFlexMessage(InquiryJson)
}

func iconCmd(_ Clients, _ *Environments, lb *LINEBot) {
	lb.SendFlexMessage(IconJson)
}

func infoCmd(c Clients, _ *Environments, lb *LINEBot) {
	go c.Postgres.InsertUserAction(lb.senderUid, UserActionInfo)
	lb.SendFlexMessage(InfoJson)
}

func loadFlexMessages(filename string, altText string) FlexMessages {
	json := richmenu.LoadJSON(filename)
	return CreateFlexMessage(json, altText)
}

func rakutanCmd(c Clients, env *Environments, lb *LINEBot) {
	go c.Postgres.InsertUserAction(lb.senderUid, UserActionRakutan)
	status, ok := GetRakutanInfo(c, env, lb.senderUid, Omikuji, Rakutan)
	if ok {
		flexMessages := CreateRakutanDetail(status.Result[0], env, Rakutan)
		lb.SendFlexMessage(flexMessages)
	} else {
		lb.SendTextMessage(status.Err)
	}
}

func onitanCmd(c Clients, env *Environments, lb *LINEBot) {
	go c.Postgres.InsertUserAction(lb.senderUid, UserActionOnitan)
	status, ok := GetRakutanInfo(c, env, lb.senderUid, Omikuji, Onitan)
	if ok {
		flexMessages := CreateRakutanDetail(status.Result[0], env, Rakutan)
		lb.SendFlexMessage(flexMessages)
	} else {
		lb.SendTextMessage(status.Err)
	}
}

func getFavoritesCmd(c Clients, env *Environments, lb *LINEBot) {
	go c.Postgres.InsertUserAction(lb.senderUid, UserActionGetFav)
	queryStatus, ok := c.Postgres.GetFavorites(lb.senderUid)
	if !ok {
		lb.SendTextMessage(queryStatus.Err)
		return
	}
	if len(queryStatus.Result) == 0 {
		lb.SendTextMessage(SuccessNoFavorites)
		return
	}
	flexMessages := CreateFavorites(queryStatus.Result)
	lb.SendFlexMessage(flexMessages)
}

func verificationCmd(c Clients, env *Environments, lb *LINEBot) {
	verified, err := c.Postgres.IsVerified(lb.senderUid)
	if err != nil {
		lb.SendTextMessage(ErrorMessageCheckVerificateError)
		return
	}
	var flexMessages FlexMessages
	if verified {
		flexMessages = VerifiedJson
	} else {
		flexMessages = VerificationJson
	}
	lb.SendFlexMessage(flexMessages)
}

func myUIDCmd(_ Clients, _ *Environments, lb *LINEBot) {
	lb.SendTextMessage(lb.senderUid)
}

func SendVerificationCmd(c Clients, env *Environments, lb *LINEBot, email string) {
	uuidObj, _ := uuid.NewUUID()
	data := []byte(lb.senderUid)
	code := uuid.NewSHA1(uuidObj, data).String()
	err := c.Postgres.InsertVerificationToken(lb.senderUid, code)
	if err != nil {
		lb.SendTextMessage(ErrorMessageInsertVerificateError)
		return
	}
	err = SendVerification(env, email, code, lb.senderUid)
	if err != nil {
		lb.SendTextMessage(ErrorMessageVerificationTokenSendError)
		return
	}
	lb.SendTextMessage(SuccessVericationTokenSent)
}
