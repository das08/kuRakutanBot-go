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

func main() {
	env := module.LoadEnv(true)
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!!")
	})

	router.GET("/verification", func(c *gin.Context) {
		uid := c.Query("uid")
		code := c.Query("code")
		postgres := module.CreatePostgresClient(&env)
		defer postgres.Client.Close()

		ok, err := postgres.CheckVerificationToken(uid, code)
		if err != nil {
			c.String(http.StatusOK, module.ErrorMessageDatabaseError)
			return
		}
		if ok {
			postgres.InsertUserAction(uid, module.UserActionVerify)
			err = postgres.UpdateUserVerification(uid)
			if err != nil {
				c.String(http.StatusOK, module.ErrorMessageDatabaseError)
				return
			}
			c.String(http.StatusOK, module.SuccessVerified)
		} else {
			c.String(http.StatusOK, module.ErrorMessageVerificationFailed)
		}
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

		postgres := module.CreatePostgresClient(&env)
		defer postgres.Client.Close()

		redis := module.CreateRedisClient()
		clients := module.Clients{Postgres: postgres, Redis: redis}

		for _, event := range events {
			postgres.IsRegistered(event.Source.UserID)
			switch event.Type {
			case linebot.EventTypeMessage:
				uid := event.Source.UserID
				lb.SetReplyToken(event.ReplyToken)
				lb.SetSenderUid(&env, uid)
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					messageText := strings.TrimSpace(message.Text)

					// コマンドが送られてきた場合
					isCommand, function := module.IsCommand(messageText)
					if isCommand {
						log.Printf("[Bot] Command: %s", messageText)
						function(clients, &env, lb)
						break
					}

					// 認証用のメールアドレスが送られてきた場合
					if module.IsStudentAddress(messageText) {
						verified, err := clients.Postgres.IsVerified(uid)
						if err != nil {
							lb.SendTextMessage(module.ErrorMessageCheckVerificateError)
						}
						if verified {
							lb.SendTextMessage(module.SuccessAlreadyVerified)
						} else {
							postgres.InsertUserAction(uid, module.UserActionEmail)
							log.Printf("[Bot] Sent verification")
							module.SendVerificationCmd(clients, &env, lb, messageText)
						}
						break
					}

					// その他講義名が送られてきた場合
					postgres.InsertUserAction(uid, module.UserActionSearch)
					searchStatus, ok := searchRakutan(clients, &env, uid, messageText)
					log.Printf("[Bot] Search: %s", messageText)
					if !ok {
						lb.SendTextMessage(searchStatus.Err)
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
						postgres.InsertUserAction(uid, module.UserActionSetFav)
						message, _ := postgres.ToggleFavorite(uid, id)
						lb.SendTextMessage(message)
					case module.Del:
						// TODO: validate
						postgres.InsertUserAction(uid, module.UserActionUnsetFav)
						message, _ := postgres.UnsetFavorite(uid, id)
						lb.SendTextMessage(message)
					}
				}
			}
		}
	})

	err := router.Run(":" + env.AppPort)
	if err != nil {
		fmt.Println("Error: creating router failed.")
		return
	}
}

func searchRakutan(c module.Clients, env *module.Environments, uid string, searchText string) (module.ExecStatus[module.FlexMessages], bool) {
	var ok, searchSuccess bool
	var status module.ExecStatus[module.RakutanInfos]
	var searchStatus module.ExecStatus[module.FlexMessages]

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
			searchStatus.Err = fmt.Sprintf(module.ErrorMessageRakutanNotFound, searchText)
		case recordCount == 1:
			favEntry, ok := c.Postgres.GetFavoriteByID(uid, rakutanInfos[0].ID)
			if ok && len(favEntry.Result) == 1 {
				rakutanInfos[0].IsFavorite = true
			}
			if !ok {
				searchStatus.Err = module.ErrorMessageGetFavError
			} else {
				searchStatus.Result = module.CreateRakutanDetail(rakutanInfos[0], env, module.Normal)
				searchSuccess = true
			}
		case recordCount <= 5*module.MaxResultsPerPage:
			searchStatus.Result = module.CreateSearchResult(searchText, rakutanInfos)
			searchSuccess = true
		default:
			searchStatus.Err = fmt.Sprintf(module.ErrorMessageTooManyRakutan, searchText, recordCount)
		}
	} else {
		searchStatus.Err = module.ErrorMessageDatabaseError
	}
	return searchStatus, searchSuccess
}
