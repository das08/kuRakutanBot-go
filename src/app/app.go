package main

import (
	"fmt"
	rakutan "github.com/das08/kuRakutanBot-go/models/rakutan"
	"github.com/das08/kuRakutanBot-go/module"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"net/http"
	"strings"
)

func main() {
	env := module.LoadEnv(true)
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!!")
	})

	router.GET("/verification", func(c *gin.Context) {
		code := c.Query("code")
		mongoDB := module.CreateDBClient(&env)
		defer mongoDB.Cancel()
		defer func() {
			//log.Println("[DB] Closed")
			if err := mongoDB.Client.Disconnect(mongoDB.Ctx); err != nil {
				panic(err)
			}
		}()
		clients := module.Clients{Mongo: mongoDB}
		res := module.CheckVerification(clients, &env, code)
		c.String(http.StatusOK, res.Message)
	})

	router.POST("/callback", func(c *gin.Context) {
		lb := module.CreateLINEBotClient(&env, c)
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
			//log.Println("[DB] Closed")
			if err := mongoDB.Client.Disconnect(mongoDB.Ctx); err != nil {
				panic(err)
			}
		}()
		redis := module.CreateRedisClient()
		clients := module.Clients{Mongo: mongoDB, Redis: redis}

		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeMessage:
				uid := event.Source.UserID
				lb.SetReplyToken(event.ReplyToken)
				lb.SetSenderUid(&env, uid)
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					messageText := strings.TrimSpace(message.Text)
					module.CountMessage(clients, &env, uid)

					// コマンドが送られてきた場合
					isCommand, function := module.IsCommand(messageText)
					if isCommand {
						log.Printf("[Bot] Command: %s", messageText)
						function(clients, &env, lb)
						break
					}

					// 認証用のメールアドレスが送られてきた場合
					if module.IsStudentAddress(messageText) {
						if module.IsVerified(clients, &env, uid) {
							lb.SendTextMessage(module.ReplyText{Status: module.KRBSuccess, Text: "すでに認証済みです。"})
						} else {
							log.Printf("[Bot] Sent verification")
							module.SendVerificationCmd(clients, &env, lb, messageText)
						}
						break
					}

					// その他講義名が送られてきた場合
					status, flexMessages := searchRakutan(clients, &env, uid, messageText)
					log.Printf("[Bot] Search: %s", messageText)
					if status.Success {
						lb.SendFlexMessage(flexMessages)
					} else {
						lb.SendTextMessage(module.ReplyText{Status: status.Status, Text: status.Message})
					}
				}
			case linebot.EventTypePostback:
				uid := event.Source.UserID
				lb.SetReplyToken(event.ReplyToken)
				lb.SetSenderUid(&env, uid)

				data := event.Postback.Data
				fmt.Println("pbdata: ", data)
				success, params := module.ParsePBParam(data)
				if success {
					fmt.Println("Params: ", params)
					switch params.Type {
					case module.Fav:
						insertStatus := module.InsertFavorite(mongoDB, &env, module.PostbackEntry{Uid: uid, Param: params})
						lb.SendTextMessage(module.ReplyText{Status: insertStatus.Status, Text: insertStatus.Message})
					case module.Del:
						deleteStatus := module.DeleteFavorite(mongoDB, &env, module.PostbackEntry{Uid: uid, Param: params})
						lb.SendTextMessage(module.ReplyText{Status: deleteStatus.Status, Text: deleteStatus.Message})
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

func searchRakutan(c module.Clients, env *module.Environments, uid string, searchText string) (module.QueryStatus, []module.FlexMessage) {
	var searchStatus module.QueryStatus
	var flexMessages []module.FlexMessage
	var queryStatus module.QueryStatus
	var result []rakutan.RakutanInfo

	isLectureNumber, lectureID := module.IsLectureID(searchText)
	if isLectureNumber {
		queryStatus, result = module.GetRakutanInfo(c, env, uid, module.ID, lectureID)
	} else {
		queryStatus, result = module.GetRakutanInfo(c, env, uid, module.Name, searchText)
	}

	if queryStatus.Success {
		recordCount := len(result)
		switch {
		case recordCount == 0:
			searchStatus.Success = false
			searchStatus.Message = fmt.Sprintf("「%s」は見つかりませんでした。\n【検索のヒント】%%を頭につけて検索すると部分一致検索ができます。ex.)「%%地理学」", searchText)
		case recordCount == 1:
			flexMessages = module.CreateRakutanDetail(result[0], module.Normal)
			searchStatus.Success = true
		case recordCount <= 5*module.MaxResultsPerPage:
			flexMessages = module.CreateSearchResult(searchText, result)
			searchStatus.Success = true
		default:
			searchStatus.Success = false
			searchStatus.Message = fmt.Sprintf("「%s」は%d件あります。検索条件を絞ってください。", searchText, recordCount)
		}
	} else {
		searchStatus.Success = false
		searchStatus.Message = "エラーが発生しました。"
		searchStatus.Status = queryStatus.Status
	}
	return searchStatus, flexMessages
}
