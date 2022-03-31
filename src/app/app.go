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

type Clients struct {
	Mongo *module.MongoDB
	Redis *module.Redis
}

func main() {
	env := module.LoadEnv(true)
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!!")
	})

	router.POST("/callback", func(c *gin.Context) {
		lb := module.CreateLINEBotClient(&env)
		events, err := lb.Bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.Writer.WriteHeader(400)
			} else {
				c.Writer.WriteHeader(500)
			}
			return
		}

		mongoDB := module.CreateDBClient(&env)
		defer mongoDB.Cancel()
		defer func() {
			fmt.Println("connection closed")
			if err := mongoDB.Client.Disconnect(mongoDB.Ctx); err != nil {
				panic(err)
			}
		}()
		redis := module.CreateRedisClient()
		clients := Clients{Mongo: mongoDB, Redis: redis}

		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeMessage:
				uid := event.Source.UserID
				lb.SetReplyToken(event.ReplyToken)
				lb.SetSenderUid(uid)
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					messageText := strings.TrimSpace(message.Text)
					module.CountMessage(clients.Mongo, &env, uid)

					isCommand, function := module.IsCommand(messageText)
					if isCommand {
						function(mongoDB, redis, &env, lb)
						break
					}

					success, flexMessages := searchRakutan(mongoDB, &env, messageText)
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
						insertStatus := module.InsertFavorite(mongoDB, &env, module.PostbackEntry{Uid: uid, Param: params})
						lb.SendTextMessage(insertStatus.Message)
					case module.Del:
						deleteStatus := module.DeleteFavorite(mongoDB, &env, module.PostbackEntry{Uid: uid, Param: params})
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

func searchRakutan(m *module.MongoDB, env *module.Environments, searchText string) (bool, []module.FlexMessage) {
	success := false
	var flexMessages []module.FlexMessage
	var queryStatus module.QueryStatus
	var result []rakutan.RakutanInfo

	isLectureNumber, lectureID := module.IsLectureID(searchText)
	if isLectureNumber {
		queryStatus, result = module.GetRakutanInfo(m, env, module.ID, lectureID)
	} else {
		queryStatus, result = module.GetRakutanInfo(m, env, module.Name, searchText)
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
