package main

import (
	"fmt"
	models "github.com/das08/kuRakutanBot-go/models/rakutan"
	status "github.com/das08/kuRakutanBot-go/models/status"
	"github.com/das08/kuRakutanBot-go/module"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"net/http"
	"strings"
)

func main() {
	env := module.LoadEnv(true)
	lb := module.CreateLINEBotClient(&env)
	//r := module.LoadRakutanDetail()
	//rr, _ := r.Marshal()
	// s, _ := r.Marshal()
	// fmt.Println(fmt.Sprintf("%s", s))
	// module.CreateDBClient(&env)

	//_, result := module.FindByLectureID(&env, mongo, 12156)
	//fmt.Printf("status: %v, result: %#v", status, result)

	//status, result := module.FindByLectureName(&env, mongo, "中国語")
	//fmt.Printf("status: %v, result: %#v", status, result)

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
			if event.Type == linebot.EventTypeMessage {
				lb.SetReplyToken(event.ReplyToken)
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					messageText := strings.TrimSpace(message.Text)

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

			}
		}
	})

	err := router.Run(":" + env.APP_PORT)
	if err != nil {
		fmt.Println("Error: creating router failed.")
		return
	}
}

//type FindByMethod int
//
//const (
//	Name FindByMethod = iota
//	ID
//	Omikuji
//)

func searchRakutan(env *module.Environments, searchText string) (bool, []module.FlexMessage) {
	success := false
	var flexMessages []module.FlexMessage
	var queryStatus status.QueryStatus
	var result []models.RakutanInfo

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
			flexMessages = module.CreateRakutanDetail(result[0])
			success = true
		default:
			flexMessages = module.CreateSearchResult(searchText, result)
			success = true
		}
	}
	return success, flexMessages
}

//func getRakutanInfo(env *module.Environments, method FindByMethod, value interface{}) (status.QueryStatus, []models.RakutanInfo) {
//	mongoDB := module.CreateDBClient(env)
//	defer mongoDB.Cancel()
//	defer mongoDB.Client.Disconnect(mongoDB.Ctx)
//	var queryStatus status.QueryStatus
//	var result []models.RakutanInfo
//
//	switch method {
//	case ID:
//		queryStatus, result = module.FindByLectureID(env, mongoDB, value.(int))
//	case Name:
//		queryStatus, result = module.FindByLectureName(env, mongoDB, value.(string))
//	case Omikuji:
//		queryStatus, result = module.FindByOmikuji(env, mongoDB, value.(string))
//	}
//
//	return queryStatus, result
//}
