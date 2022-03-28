package module

import "strings"

type Command struct {
	Keyword      string
	SendFunction func(env *Environments, lb *LINEBot)
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
}

func IsCommand(messageText string) (bool, func(env *Environments, lb *LINEBot)) {
	isCommand := false
	var function func(env *Environments, lb *LINEBot)
	for _, cmd := range Commands {
		// Case-insensitive
		if strings.EqualFold(cmd.Keyword, messageText) {
			isCommand = true
			function = cmd.SendFunction
		}
	}
	return isCommand, function
}

func helpCmd(_ *Environments, lb *LINEBot) {
	help := LoadHelp()
	helpJson, _ := help.Marshal()
	flexMessages := CreateFlexMessage(helpJson, "ヘルプ")
	lb.SendFlexMessage(flexMessages)
}

func judgeDetailCmd(_ *Environments, lb *LINEBot) {
	judgeDetail := LoadJudgeDetail()
	judgeDetailJson, _ := judgeDetail.Marshal()
	flexMessages := CreateFlexMessage(judgeDetailJson, "らくたん判定の詳細")
	lb.SendFlexMessage(flexMessages)
}

func rakutanCmd(env *Environments, lb *LINEBot) {
	queryStatus, result := GetRakutanInfo(env, Omikuji, "rakutan")
	if queryStatus.Success {
		flexMessages := CreateRakutanDetail(result[0])
		lb.SendFlexMessage(flexMessages)
	}
}

func onitanCmd(env *Environments, lb *LINEBot) {
	queryStatus, result := GetRakutanInfo(env, Omikuji, "onitan")
	if queryStatus.Success {
		flexMessages := CreateRakutanDetail(result[0])
		lb.SendFlexMessage(flexMessages)
	}
}
