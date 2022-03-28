package module

import (
	"fmt"
	models "github.com/das08/kuRakutanBot-go/models/rakutan"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"strconv"
)

type FlexMessage struct {
	FlexContainer linebot.FlexContainer
	AltText       string
}

func CreateRakutanDetail(info models.RakutanInfo) []FlexMessage {
	rakutanDetail := LoadRakutanDetail()
	rakutanDetail.Header.Contents[0].Contents[1].Text = toPtr("Search ID:#" + toStr(info.ID))
	rakutanDetail.Header.Contents[1].Text = &info.LectureName             // Lecture name
	rakutanDetail.Header.Contents[3].Contents[1].Text = &info.FacultyName // Faculty name
	rakutanDetail.Header.Contents[4].Contents[1].Text = toPtr("---")      // Group
	rakutanDetail.Header.Contents[4].Contents[3].Text = toPtr("---")      // Credits

	// 単位取得率
	for i, v := range info.Detail {
		rakutanDetail.Body.Contents[0].Contents[i+1].Contents[0].Text = toStr(v.Year) + "年度"
		rakutanDetail.Body.Contents[0].Contents[i+1].Contents[1].Text = getRakutanPercent(v.Accepted, v.Total)
	}
	rakutanJudge := getRakutanJudge(info.Detail)
	rakutanDetail.Body.Contents[0].Contents[5].Contents[1].Text = rakutanJudge.rank
	rakutanDetail.Body.Contents[0].Contents[5].Contents[1].Color = rakutanJudge.color

	flex, err := rakutanDetail.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	flexContainer, _ := linebot.UnmarshalFlexMessageJSON(flex)
	altText := "「" + info.LectureName + "」のらくたん情報"

	return []FlexMessage{{FlexContainer: flexContainer, AltText: altText}}
}

func getRakutanPercent(accept int, total int) string {
	breakdown := "(" + toStr(accept) + "/" + toStr(total) + ")"
	if total == 0 {
		return "---% " + breakdown
	} else {
		return fmt.Sprintf("%.1f%% ", getPercentage(accept, total)) + breakdown
	}
}

type RakutanJudge struct {
	percentBound float32
	rank         string
	color        string
}

var judgeList = [9]RakutanJudge{
	{percentBound: 90, rank: "SSS", color: "#c3c45b"},
	{percentBound: 85, rank: "SS", color: "#c3c45b"},
	{percentBound: 80, rank: "S", color: "#c3c45b"},
	{percentBound: 75, rank: "A", color: "#cf2904"},
	{percentBound: 70, rank: "B", color: "#098ae0"},
	{percentBound: 60, rank: "C", color: "#f48a1c"},
	{percentBound: 50, rank: "D", color: "#8a30c9"},
	{percentBound: 0, rank: "F", color: "#837b8a"},
	{percentBound: -1, rank: "---", color: "#837b8a"},
}

func getRakutanJudge(rds models.RakutanDetails) RakutanJudge {
	accept, total := rds.GetLatestDetail()
	if total == 0 {
		return judgeList[8]
	}

	percentage := getPercentage(accept, total)
	var res = judgeList[8]
	for i, j := range judgeList {
		if percentage >= j.percentBound {
			res = judgeList[i]
			break
		}
	}
	return res
}

func getPercentage(accept int, total int) float32 {
	return float32(100*accept) / float32(total)
}

func toStr(i int) string {
	return strconv.Itoa(i)
}
func toPtr(s string) *string {
	return &s
}
