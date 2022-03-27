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
	mongo := module.CreateDBClient(&env)
	defer mongo.Cancel()
	defer mongo.Client.Disconnect(mongo.Ctx)
	lb := module.CreateLINEBotClient(&env)
	r := module.LoadRakutanDetail()
	rr, _ := r.Marshal()
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
					//lb.SendTextMessage(message.Text)
					lb.SendFlexMessage(rr, message.Text)
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
