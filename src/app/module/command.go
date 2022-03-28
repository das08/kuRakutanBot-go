package module

import "strings"

type Command struct {
	Keyword  string
	Function func(lb *LINEBot)
}

var Commands = [...]Command{
	{Keyword: "help", Function: helpCmd},
	{Keyword: "へるぷ", Function: helpCmd},
	{Keyword: "ヘルプ", Function: helpCmd},
	{Keyword: "はんてい", Function: judgeDetailCmd},
	{Keyword: "判定", Function: judgeDetailCmd},
	{Keyword: "詳細", Function: judgeDetailCmd},
	{Keyword: "判定詳細", Function: judgeDetailCmd},
	{Keyword: "楽単詳細", Function: judgeDetailCmd},
}

func IsCommand(messageText string) (bool, func(lb *LINEBot)) {
	isCommand := false
	var function func(lb *LINEBot)
	for _, cmd := range Commands {
		// Case-insensitive
		if strings.EqualFold(cmd.Keyword, messageText) {
			isCommand = true
			function = cmd.Function
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
