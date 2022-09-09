package module

import (
	"fmt"
	"github.com/das08/kuRakutanBot-go/models/richmenu"
	"github.com/jackc/pgtype"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"math"
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

type OmikujiType string

const (
	Normal  OmikujiType = "normal"
	Rakutan OmikujiType = "rakutan"
	Onitan  OmikujiType = "onitan"
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

func CreateRakutanDetail(info RakutanInfo2, e *Environments, o OmikujiType) []FlexMessage {
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
	rakutanPercents := getRakutanPercent(info.Passed, info.Register)
	for i := range info.Register.Elements {
		rakutanDetail.Body.Contents[0].Contents[i+1].Contents[0].Text = fmt.Sprintf("%d年度", e.YEAR-i)
		rakutanDetail.Body.Contents[0].Contents[i+1].Contents[1].Text = rakutanPercents[i]
	}
	rakutanJudge := getRakutanJudge(info)
	offset := len(info.Register.Elements)
	rakutanDetail.Body.Contents[0].Contents[offset+2].Contents[1].Text = rakutanJudge.rank
	rakutanDetail.Body.Contents[0].Contents[offset+2].Contents[1].Color = rakutanJudge.color

	// 過去問リンク
	// TODO: なんとかする
	//if info.IsVerified {
	//	_, err := url.ParseRequestURI(info.URL)
	//	switch {
	//	case info.URL != "" && err == nil:
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Text = "○"
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Color = "#0fd142"
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Text = "リンク"
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Color = "#4c7cf5"
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Action.URI = &info.URL
	//	case info.KUWikiErr != "":
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Text = "×"
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Color = "#ef1d2f"
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Text = info.KUWikiErr
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Action = &richmenu.URIAction{Type: "uri", Label: "action", URI: strToPtr("https://www.kuwiki.net/upload-exams")}
	//	default:
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Text = "×"
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Color = "#ef1d2f"
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Text = "提供する"
	//		rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Action = &richmenu.URIAction{Type: "uri", Label: "action", URI: strToPtr("https://www.kuwiki.net/upload-exams")}
	//	}
	//} else {
	//	rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Flex = intToPtr(0)
	//	rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Text = "△"
	//	rakutanDetail.Body.Contents[0].Contents[6].Contents[1].Color = "#ffb101"
	//	rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Flex = intToPtr(7)
	//	rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Text = "ユーザー認証が必要です"
	//	rakutanDetail.Body.Contents[0].Contents[6].Contents[2].Action = &richmenu.URIAction{Type: "message", Label: "action", Text: strToPtr("ユーザ認証")}
	//}

	flex, err := rakutanDetail.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	flexContainer, _ := linebot.UnmarshalFlexMessageJSON(flex)
	altText := fmt.Sprintf("「%s」のらくたん情報", info.LectureName)

	return []FlexMessage{{FlexContainer: flexContainer, AltText: altText}}
}

func CreateSearchResult(searchText string, infos []RakutanInfo2) []FlexMessage {
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

func CreateFavorites2(r []RakutanInfo2) []FlexMessage {
	var messages []FlexMessage
	favorites := LoadFavorites()

	pageCount := 0
	maxPageCount := len(r)/MaxResultsPerPage + 1

	for pageCount = 1; pageCount <= maxPageCount; pageCount++ {
		altText := fmt.Sprintf("お気に入りリスト(%d/%d)", pageCount, maxPageCount)
		// Set body text
		favorites.Body.Contents[0].Contents = getFavoriteList2(r, pageCount)
		flexContainer := toFlexContainer(&favorites)
		messages = append(messages, FlexMessage{FlexContainer: flexContainer, AltText: altText})
	}
	return messages
}

func CreateFlexMessage(flex []byte, altText string) []FlexMessage {
	flexContainer, _ := linebot.UnmarshalFlexMessageJSON(flex)
	return []FlexMessage{{FlexContainer: flexContainer, AltText: altText}}
}

func getLectureList(infos []RakutanInfo2, pageCount int) []richmenu.PurpleContent {
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

func getFavoriteList2(r []RakutanInfo2, pageCount int) []richmenu.FavoritesBodyContents {
	favorites := LoadFavorites()
	var favoriteList []richmenu.FavoritesBodyContents
	favorite := favorites.Body.Contents[0].Contents[0]

	offset := (pageCount - 1) * MaxResultsPerPage
	for i := offset; i < int(math.Min(float64(len(r)), float64(MaxResultsPerPage+offset))); i++ {
		tmp := favorite.DeepCopy()
		tmp.Contents[0].Text = r[i].LectureName
		tmp.Contents[1].Action.Text = strToPtr(fmt.Sprintf("#%d", r[i].ID))
		tmp.Contents[2].Action.Data = strToPtr(fmt.Sprintf("type=del&id=%d&lecname=%s", r[i].ID, r[i].LectureName))

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

func getRakutanPercent(passed pgtype.Int2Array, register pgtype.Int2Array) []string {
	var rakutanPercent []string
	for i := 0; i < len(passed.Elements); i++ {
		p := int(passed.Elements[i].Int)
		r := int(register.Elements[i].Int)
		breakdown := fmt.Sprintf("(%d/%d)", p, r)
		if register.Elements[i].Status == pgtype.Null {
			rakutanPercent = append(rakutanPercent, "---% "+breakdown)
		} else {
			rakutanPercent = append(rakutanPercent, fmt.Sprintf("%.1f%% ", getPercentage(p, r))+breakdown)
		}
	}
	return rakutanPercent
}

func getRakutanJudge(r RakutanInfo2) RakutanJudge {
	accept, total := r.GetLatestDetail()
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
