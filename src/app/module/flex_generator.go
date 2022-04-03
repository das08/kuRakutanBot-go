package module

import (
	"fmt"
	models "github.com/das08/kuRakutanBot-go/models/rakutan"
	"github.com/das08/kuRakutanBot-go/models/richmenu"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"math"
	"net/url"
	"strconv"
)

type FlexMessage struct {
	FlexContainer linebot.FlexContainer
	AltText       string
}

type RakutanJudge struct {
	percentBound float32
	rank         string
	color        string
}

type OmikujiType int

const (
	Normal OmikujiType = iota
	Rakutan
	Onitan
)

type OmikujiText struct {
	Text  string
	Color string
}

var MaxResultsPerPage = 25
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
var facultyAbbr = map[string]string{
	"文学部": "文", "教育学部": "教", "法学部": "法", "経済学部": "経", "理学部": "理", "医学部": "医医",
	"医学部（人間健康科学科）": "人健", "医学部(人間健康科学科)": "人健",
	"薬学部": "薬", "工学部": "工", "農学部": "農", "総合人間学部": "総人", "国際高等教育院": "般教",
}

var omikujiType = map[OmikujiType]OmikujiText{
	Rakutan: {Text: "楽単おみくじ結果", Color: "#ff7e41"},
	Onitan:  {Text: "鬼単おみくじ結果", Color: "#6d7bff"},
}

func CreateRakutanDetail(info models.RakutanInfo, o OmikujiType) []FlexMessage {
	rakutanDetail := LoadRakutanDetail()
	rakutanDetail.Header.Contents[0].Contents[1].Text = strToPtr("Search ID:#" + toStr(info.ID))
	rakutanDetail.Header.Contents[1].Text = &info.LectureName             // Lecture name
	rakutanDetail.Header.Contents[3].Contents[1].Text = &info.FacultyName // Faculty name
	rakutanDetail.Header.Contents[4].Contents[1].Text = strToPtr("---")   // Group
	rakutanDetail.Header.Contents[4].Contents[3].Text = strToPtr("---")   // Credits

	if info.IsFavorite {
		rakutanDetail.Header.Contents[0].Contents[0].URL = strToPtr("https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png")
	}

	if o != Normal {
		rakutanDetail.Header.Contents[0].Contents[1].Text = strToPtr(omikujiType[o].Text)
		rakutanDetail.Header.Contents[0].Contents[1].Color = strToPtr(omikujiType[o].Color)
	}

	// Postbackパラメータ
	rakutanDetail.Header.Contents[0].Contents[0].Action.Data = fmt.Sprintf("type=fav&id=%d&lecname=%s", info.ID, info.LectureName)

	// 単位取得率
	for i, v := range info.Detail {
		rakutanDetail.Body.Contents[0].Contents[i+1].Contents[0].Text = fmt.Sprintf("%d年度", v.Year)
		rakutanDetail.Body.Contents[0].Contents[i+1].Contents[1].Text = getRakutanPercent(v.Accepted, v.Total)
	}
	rakutanJudge := getRakutanJudge(info.Detail)
	rakutanDetail.Body.Contents[0].Contents[5].Contents[1].Text = rakutanJudge.rank
	rakutanDetail.Body.Contents[0].Contents[5].Contents[1].Color = rakutanJudge.color

	// 過去問リンク
	if info.IsVerified {
		_, err := url.ParseRequestURI(info.URL)
		switch {
		case info.URL != "" && err == nil:
			rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Text = "○"
			rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Color = "#0fd142"
			rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Text = "リンク"
			rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Color = "#4c7cf5"
			rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Action.URI = &info.URL
		case info.KUWikiErr != "":
			rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Text = "×"
			rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Color = "#ef1d2f"
			rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Text = info.KUWikiErr
			rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Action = &richmenu.URIAction{Type: "uri", Label: "action", URI: strToPtr("https://www.kuwiki.net/upload-exams")}
		default:
			rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Text = "×"
			rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Color = "#ef1d2f"
			rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Text = "提供する"
			rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Action = &richmenu.URIAction{Type: "uri", Label: "action", URI: strToPtr("https://www.kuwiki.net/upload-exams")}
		}
	} else {
		rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Flex = intToPtr(0)
		rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Text = "△"
		rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Color = "#ffb101"
		rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Flex = intToPtr(7)
		rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Text = "ユーザー認証が必要です"
		rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Action = &richmenu.URIAction{Type: "message", Label: "action", Text: strToPtr("ユーザ認証")}
	}

	flex, err := rakutanDetail.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	flexContainer, _ := linebot.UnmarshalFlexMessageJSON(flex)
	altText := fmt.Sprintf("「%s」のらくたん情報", info.LectureName)

	return []FlexMessage{{FlexContainer: flexContainer, AltText: altText}}
}

func CreateSearchResult(searchText string, infos []models.RakutanInfo) []FlexMessage {
	var messages []FlexMessage
	searchResult := LoadSearchResult()
	searchResultMore := LoadSearchResultMore()

	pageCount := 0
	maxPageCount := len(infos)/MaxResultsPerPage + 1

	for pageCount = 1; pageCount <= maxPageCount; pageCount++ {
		altText := fmt.Sprintf("「%s」の検索結果(%d/%d)", searchText, pageCount, maxPageCount)
		switch {
		case pageCount == 1:
			// Set header text
			searchResult.Header.Contents[0].Text = strToPtr(altText)
			searchResult.Header.Contents[1].Text = strToPtr(fmt.Sprintf("%d 件の候補が見つかりました。目的の講義を選択してください。", len(infos)))

			// Set body text
			searchResult.Body.Contents[1].Contents = getLectureList(infos, pageCount)
			flexContainer := toFlexContainer(&searchResult)
			messages = append(messages, FlexMessage{FlexContainer: flexContainer, AltText: altText})
		case pageCount >= 2:
			// Set header text
			searchResultMore.Header.Contents[0].Text = strToPtr(altText)

			// Set body text
			searchResultMore.Body.Contents[1].Contents = getLectureList(infos, pageCount)
			flexContainer := toFlexContainer(&searchResultMore)
			messages = append(messages, FlexMessage{FlexContainer: flexContainer, AltText: altText})
		}
	}
	return messages
}

func CreateFavorites(favs []models.Favorite) []FlexMessage {
	var messages []FlexMessage
	favorites := LoadFavorites()

	pageCount := 0
	maxPageCount := len(favs)/MaxResultsPerPage + 1

	for pageCount = 1; pageCount <= maxPageCount; pageCount++ {
		altText := fmt.Sprintf("お気に入りリスト(%d/%d)", pageCount, maxPageCount)
		// Set body text
		favorites.Body.Contents[0].Contents = getFavoriteList(favs, pageCount)
		flexContainer := toFlexContainer(&favorites)
		messages = append(messages, FlexMessage{FlexContainer: flexContainer, AltText: altText})
	}
	return messages
}

func CreateFlexMessage(flex []byte, altText string) []FlexMessage {
	flexContainer, _ := linebot.UnmarshalFlexMessageJSON(flex)
	return []FlexMessage{{FlexContainer: flexContainer, AltText: altText}}
}

func getLectureList(infos []models.RakutanInfo, pageCount int) []richmenu.PurpleContent {
	searchResult := LoadSearchResult()
	var lectureList []richmenu.PurpleContent
	lecture := searchResult.Body.Contents[1].Contents[0]

	offset := (pageCount - 1) * MaxResultsPerPage
	for i := offset; i < int(math.Min(float64(len(infos)), float64(MaxResultsPerPage+offset))); i++ {
		tmp := lecture.DeepCopy()
		tmp.Contents[1].Text = infos[i].LectureName
		tmp.Contents[2].Action.Text = fmt.Sprintf("#%d", infos[i].ID)

		abbr, ok := facultyAbbr[infos[i].FacultyName]
		if ok {
			tmp.Contents[0].Text = abbr
		}

		lectureList = append(lectureList, tmp)
	}

	return lectureList
}

func getFavoriteList(favs []models.Favorite, pageCount int) []richmenu.FavoritesBodyContents {
	favorites := LoadFavorites()
	var favoriteList []richmenu.FavoritesBodyContents
	favorite := favorites.Body.Contents[0].Contents[0]

	offset := (pageCount - 1) * MaxResultsPerPage
	for i := offset; i < int(math.Min(float64(len(favs)), float64(MaxResultsPerPage+offset))); i++ {
		tmp := favorite.DeepCopy()
		tmp.Contents[0].Text = favs[i].LectureName
		tmp.Contents[1].Action.Text = strToPtr(fmt.Sprintf("#%d", favs[i].ID))
		tmp.Contents[2].Action.Data = strToPtr(fmt.Sprintf("type=del&id=%d&lecname=%s", favs[i].ID, favs[i].LectureName))

		favoriteList = append(favoriteList, tmp)
	}

	return favoriteList
}

func toFlexContainer(json richmenu.Marshal) linebot.FlexContainer {
	flex, err := json.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	flexContainer, _ := linebot.UnmarshalFlexMessageJSON(flex)
	return flexContainer
}

func getRakutanPercent(accept int, total int) string {
	breakdown := "(" + toStr(accept) + "/" + toStr(total) + ")"
	if total == 0 {
		return "---% " + breakdown
	} else {
		return fmt.Sprintf("%.1f%% ", getPercentage(accept, total)) + breakdown
	}
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

func strToPtr(s string) *string {
	return &s
}

func intToPtr(i int) *int {
	return &i
}
