package main

import (
	"fmt"
	rakutan "github.com/das08/kuRakutanBot-go/models/rakutan"
	"github.com/das08/kuRakutanBot-go/module"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"net/http"
	"strings"
)

func main() {
	env := module.LoadEnv(true)
	lb := module.CreateLINEBotClient(&env)

	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!!")
	})

	router.POST("/callback", func(c *gin.Context) {
		events, err := lb.Bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.Writer.WriteHeader(400)
			} else {
				c.Writer.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeMessage:
				uid := event.Source.UserID
				lb.SetReplyToken(event.ReplyToken)
				lb.SetSenderUid(uid)
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					messageText := strings.TrimSpace(message.Text)
					module.CountMessage(&env, uid)

					isCommand, function := module.IsCommand(messageText)
					if isCommand {
						function(&env, lb)
						break
					}

					success, flexMessages := searchRakutan(&env, messageText)
					if success {
						lb.SendFlexMessage(flexMessages)
					} else {
						lb.SendTextMessage(message.Text)
					}
				}
			case linebot.EventTypePostback:
				uid := event.Source.UserID
				lb.SetReplyToken(event.ReplyToken)
				lb.SetSenderUid(uid)

				data := event.Postback.Data
				fmt.Println("pbdata: ", data)
				success, params := module.ParsePBParam(data)
				if success {
					fmt.Println("Params: ", params)
					switch params.Type {
					case module.Fav:
						insertStatus := module.InsertFavorite(&env, module.PostbackEntry{Uid: uid, Param: params})
						lb.SendTextMessage(insertStatus.Message)
					case module.Del:
						deleteStatus := module.DeleteFavorite(&env, module.PostbackEntry{Uid: uid, Param: params})
						lb.SendTextMessage(deleteStatus.Message)
					}
				}
			}
		}
	})

	err := router.Run(":" + env.APP_PORT)
	if err != nil {
		fmt.Println("Error: creating router failed.")
		return
	}
}

func searchRakutan(env *module.Environments, searchText string) (bool, []module.FlexMessage) {
	success := false
	var flexMessages []module.FlexMessage
	var queryStatus module.QueryStatus
	var result []rakutan.RakutanInfo

	isLectureNumber, lectureID := module.IsLectureID(searchText)
	if isLectureNumber {
		queryStatus, result = module.GetRakutanInfo(env, module.ID, lectureID)
	} else {
		queryStatus, result = module.GetRakutanInfo(env, module.Name, searchText)
	}

	if queryStatus.Success {
		recordCount := len(result)
		switch recordCount {
		case 0:
			break
		case 1:
			flexMessages = module.CreateRakutanDetail(result[0], module.Normal)
			success = true
		default:
			flexMessages = module.CreateSearchResult(searchText, result)
			success = true
		}
	}
	return success, flexMessages
}
