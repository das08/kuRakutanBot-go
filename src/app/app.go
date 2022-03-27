package main

import (
	"fmt"
	"github.com/das08/kuRakutanBot-go/module"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"net/http"
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
					fmt.Println(message.Text)

					success, flex := searchRakutan(&env, message.Text)
					if success {
						lb.SendFlexMessage(flex, message.Text)
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

func searchRakutan(env *module.Environments, searchText string) (bool, []byte) {
	success := false
	var flex []byte
	mongoDB := module.CreateDBClient(env)
	defer mongoDB.Cancel()
	defer mongoDB.Client.Disconnect(mongoDB.Ctx)

	status, result := module.FindByLectureName(env, mongoDB, searchText)

	if status.Success {
		recordCount := len(result)
		fmt.Println(recordCount)
		if recordCount == 1 {
			flex = module.SetRakutanData(result[0])
			success = true
		}
	}

	return success, flex
}
