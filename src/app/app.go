package main

import (
	"fmt"
	"github.com/das08/kuRakutanBot-go/module"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"net/http"
	"strings"
)

const (
	YEAR = "year"
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

		//mongoDB := module.CreateDBClient(&env)
		//defer mongoDB.Cancel()
		//defer func() {
		//	//log.Println("[DB] Closed")
		//	if err := mongoDB.Client.Disconnect(mongoDB.Ctx); err != nil {
		//		panic(err)
		//	}
		//}()

		postgres := module.CreatePostgresClient(&env)
		defer postgres.Client.Close()

		redis := module.CreateRedisClient()
		clients := module.Clients{Postgres: postgres, Redis: redis}

		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeMessage:
				uid := event.Source.UserID
				lb.SetReplyToken(event.ReplyToken)
				lb.SetSenderUid(&env, uid)
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					messageText := strings.TrimSpace(message.Text)
					//module.CountMessage(clients, &env, uid)

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
					searchStatus, ok := searchRakutan(clients, &env, uid, messageText)
					log.Printf("[Bot] Search: %s", messageText)
					if !ok {
						lb.SendTextMessage2(searchStatus.Err)
					} else {
						lb.SendFlexMessage(searchStatus.Result)
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
					id := params.ID
					switch params.Type {
					case module.Fav:
						// TODO: validate
						message, _ := postgres.SetFavorite(uid, id)
						lb.SendTextMessage2(message)
					case module.Del:
						// TODO: validate
						message, _ := postgres.UnsetFavorite(uid, id)
						lb.SendTextMessage2(message)
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

func searchRakutan(c module.Clients, env *module.Environments, uid string, searchText string) (module.QueryStatus2[[]module.FlexMessage], bool) {
	var ok, searchSuccess bool
	var status module.QueryStatus2[[]module.RakutanInfo2]
	var searchStatus module.QueryStatus2[[]module.FlexMessage]

	isLectureNumber, lectureID := module.IsLectureID(searchText)
	if isLectureNumber {
		status, ok = module.GetRakutanInfo(c, env, uid, module.ID, lectureID)
	} else {
		status, ok = module.GetRakutanInfo(c, env, uid, module.Name, searchText)
	}

	if ok {
		rakutanInfos := status.Result
		recordCount := len(rakutanInfos)
		switch {
		case recordCount == 0:
			searchStatus.Err = fmt.Sprintf("「%s」は見つかりませんでした。\n【検索のヒント】%%を頭につけて検索すると部分一致検索ができます。ex.)「%%地理学」", searchText)
			//searchSuccess = false

		case recordCount == 1:
			favEntry, ok := c.Postgres.GetFavoriteByID(uid, rakutanInfos[0].ID)
			if ok && len(favEntry.Result) == 1 {
				rakutanInfos[0].IsFavorite = true
			}
			if !ok {
				searchStatus.Err = "お気に入りの取得に失敗しました。"
				//searchSuccess = false
			} else {
				searchStatus.Result = module.CreateRakutanDetail(rakutanInfos[0], env, module.Normal)
				searchSuccess = true
			}

		case recordCount <= 5*module.MaxResultsPerPage:
			searchStatus.Result = module.CreateSearchResult(searchText, rakutanInfos)
			searchSuccess = true

		default:
			searchStatus.Err = fmt.Sprintf("「%s」は%d件あります。検索条件を絞ってください。", searchText, recordCount)
			//searchSuccess = false
		}
	} else {
		searchStatus.Err = "エラーが発生しました。"
		//searchSuccess = false
	}
	return searchStatus, searchSuccess
}
