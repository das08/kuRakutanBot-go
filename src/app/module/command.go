package module

import "strings"

type Command struct {
	Keyword      string
	DBFunction   func()
	SendFunction func(lb *LINEBot)
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
}

func IsCommand(messageText string) (bool, func(lb *LINEBot)) {
	isCommand := false
	var function func(lb *LINEBot)
	for _, cmd := range Commands {
		// Case-insensitive
		if strings.EqualFold(cmd.Keyword, messageText) {
			isCommand = true
			function = cmd.SendFunction
		}
	}
	return isCommand, function
}

func helpCmd(lb *LINEBot) {
	help := LoadHelp()
	helpJson, _ := help.Marshal()
	flexMessages := CreateFlexMessage(helpJson, "ヘルプ")
	lb.SendFlexMessage(flexMessages)
}

func judgeDetailCmd(lb *LINEBot) {
	judgeDetail := LoadJudgeDetail()
	judgeDetailJson, _ := judgeDetail.Marshal()
	flexMessages := CreateFlexMessage(judgeDetailJson, "らくたん判定の詳細")
	lb.SendFlexMessage(flexMessages)
}

//func omikujiCmd(lb *LINEBot) {
//	queryStatus, result := FindByOmikuji(env, mongoDB, "rakutan")
//}
