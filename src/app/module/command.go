package module

import "strings"

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
	{Keyword: "お問い合わせ", SendFunction: inquiryCmd},
	{Keyword: "問い合わせ", SendFunction: inquiryCmd},
	{Keyword: "京大楽単bot", SendFunction: infoCmd},
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

func infoCmd(_ Clients, _ *Environments, lb *LINEBot) {
	flexMessages := loadFlexMessages("./assets/richmenu/info.json", "京大楽単bot")
	lb.SendFlexMessage(flexMessages)
}

func loadFlexMessages(filename string, altText string) []FlexMessage {
	json := LoadJSON(filename)
	return CreateFlexMessage(json, altText)
}

func rakutanCmd(c Clients, env *Environments, lb *LINEBot) {
	queryStatus, result := GetRakutanInfo(c, env, Omikuji, "rakutan")
	countUp(env, c.Mongo, lb.senderUid, "rakutan")
	if queryStatus.Success {
		flexMessages := CreateRakutanDetail(result[0], Rakutan)
		lb.SendFlexMessage(flexMessages)
	} else {
		lb.SendTextMessage("楽単おみくじに失敗しました。")
	}
}

func onitanCmd(c Clients, env *Environments, lb *LINEBot) {
	queryStatus, result := GetRakutanInfo(c, env, Omikuji, "onitan")
	countUp(env, c.Mongo, lb.senderUid, "onitan")
	if queryStatus.Success {
		flexMessages := CreateRakutanDetail(result[0], Onitan)
		lb.SendFlexMessage(flexMessages)
	} else {
		lb.SendTextMessage("鬼単おみくじに失敗しました。")
	}
}
func getFavoritesCmd(c Clients, env *Environments, lb *LINEBot) {
	queryStatus, result := GetFavorites(c, env, lb.senderUid)
	if queryStatus.Success {
		flexMessages := CreateFavorites(result)
		lb.SendFlexMessage(flexMessages)
	} else {
		lb.SendTextMessage(queryStatus.Message)
	}
}
