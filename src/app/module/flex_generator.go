package module

import (
	"fmt"
	richmenu2 "github.com/das08/kuRakutanBot-go/assets/richmenu"
	"github.com/das08/kuRakutanBot-go/richmenu"
	"github.com/jackc/pgtype"
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

var MaxResultsPerPage = 30
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
	"文学研究科": "文研", "教育学研究科": "教研", "法学研究科": "法研", "経済学研究科": "経研", "理学研究科": "理研",
	"医学研究科": "医研", "医学研究科(人間健康科学系専攻": "医研", "薬学研究科": "薬研", "工学研究科": "工研",
	"農学研究科": "農研", "人間・環境学研究科": "人環", "エネルギー科学研究科": "エネ", "アジア・アフリカ地域研究研究科": "アア",
	"情報学研究科": "情研", "生命科学研究科": "生命", "地球環境学舎": "地環", "公共政策教育部": "公政",
	"経営管理教育部": "経営", "法学研究科(法科大学院)": "法科", "総合生存学館": "生存",
}

var omikujiType = map[OmikujiType]OmikujiText{
	Rakutan: {Text: "楽単おみくじ結果", Color: "#ff7e41"},
	Onitan:  {Text: "鬼単おみくじ結果", Color: "#6d7bff"},
}

func CreateRakutanDetail(info RakutanInfo, e *Environments, o OmikujiType) FlexMessages {
	flexContainer := richmenu2.LoadRakutanDetail()
	flexContainer.Header.Contents[0].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Text = "Search ID:#" + toStr(info.ID)
	flexContainer.Header.Contents[1].(*linebot.TextComponent).Text = info.LectureName                                     // Lecture name
	flexContainer.Header.Contents[3].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Text = info.FacultyName // Faculty name
	flexContainer.Header.Contents[4].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Text = "---"            // Group
	flexContainer.Header.Contents[4].(*linebot.BoxComponent).Contents[3].(*linebot.TextComponent).Text = "---"            // Credits

	if info.IsFavorite {
		flexContainer.Header.Contents[0].(*linebot.BoxComponent).Contents[0].(*linebot.ImageComponent).URL = "https://scdn.line-apps.com/n/channel_devcenter/img/fx/review_gold_star_28.png"
	}

	if o != Normal {
		flexContainer.Header.Contents[0].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Text = omikujiType[o].Text
		flexContainer.Header.Contents[0].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Color = omikujiType[o].Color
	}

	// Postbackパラメータ
	flexContainer.Header.Contents[0].(*linebot.BoxComponent).Contents[0].(*linebot.ImageComponent).Action.(*linebot.PostbackAction).Data = fmt.Sprintf("type=fav&id=%d&lecname=%s", info.ID, info.LectureName)

	// 単位取得率
	rakutanPercents := getRakutanPercent(info.Passed, info.Register)
	for i := range info.Register.Elements {
		flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[i+1].(*linebot.BoxComponent).Contents[0].(*linebot.TextComponent).Text = fmt.Sprintf("%d年度", e.YEAR-i)
		flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[i+1].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Text = rakutanPercents[i]
	}
	rakutanJudge := getRakutanJudge(info)
	offset := len(info.Register.Elements)
	flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+2].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Text = rakutanJudge.rank
	flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+2].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Color = rakutanJudge.color

	// 過去問リンク
	// TODO: なんとかする
	if info.IsVerified {
		_, err := url.ParseRequestURI(info.KakomonURL)
		switch {
		case info.KakomonURL != "" && err == nil:
			flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Text = "○"
			flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Color = "#0fd142"
			flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[2].(*linebot.TextComponent).Text = "リンク"
			flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[2].(*linebot.TextComponent).Color = "#4c7cf5"
			flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[2].(*linebot.TextComponent).Action.(*linebot.URIAction).URI = info.KakomonURL
		default:
			flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Text = "×"
			flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Color = "#ef1d2f"
			flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[2].(*linebot.TextComponent).Text = "提供する"
			flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[2].(*linebot.TextComponent).Action.(*linebot.URIAction).URI = "https://www.kuwiki.net/upload-exams"
		}
	} else {
		flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Flex = intToPtr(0)
		flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Text = "△"
		flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[1].(*linebot.TextComponent).Color = "#ffb101"
		flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[2].(*linebot.TextComponent).Flex = intToPtr(7)
		flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[2].(*linebot.TextComponent).Text = "ユーザー認証が必要です"
		//flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[2].(*linebot.TextComponent).Action.(*linebot.MessageAction).Text = "ユーザ認証"
		flexContainer.Body.Contents[0].(*linebot.BoxComponent).Contents[offset+3].(*linebot.BoxComponent).Contents[2].(*linebot.TextComponent).Action = &linebot.MessageAction{
			Label: "action",
			Text:  "ユーザ認証",
		}

	}

	altText := fmt.Sprintf("「%s」のらくたん情報", info.LectureName)

	return FlexMessages{{FlexContainer: flexContainer, AltText: altText}}
}

func CreateSearchResult(searchText string, infos RakutanInfos) FlexMessages {
	var messages FlexMessages
	searchResult := richmenu.SearchResultJson
	searchResultMore := richmenu.SearchResultMoreJson

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

func CreateFavorites(r RakutanInfos) FlexMessages {
	var messages FlexMessages
	favorites := richmenu.FavoritesJson

	pageCount := 0
	maxPageCount := len(r)/MaxResultsPerPage + 1

	for pageCount = 1; pageCount <= maxPageCount; pageCount++ {
		altText := fmt.Sprintf("お気に入りリスト(%d/%d)", pageCount, maxPageCount)
		// Set body text
		favorites.Body.Contents[0].Contents = getFavoriteList(r, pageCount)
		flexContainer := toFlexContainer(&favorites)
		messages = append(messages, FlexMessage{FlexContainer: flexContainer, AltText: altText})
	}
	return messages
}

func CreateFlexMessage(flex []byte, altText string) FlexMessages {
	flexContainer, err := linebot.UnmarshalFlexMessageJSON(flex)
	if err != nil {
		log.Fatal(err)
	}
	return FlexMessages{{FlexContainer: flexContainer, AltText: altText}}
}

func getLectureList(infos RakutanInfos, pageCount int) []richmenu.PurpleContent {
	searchResult := richmenu.SearchResultJson
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
		} else {
			tmp.Contents[0].Text = "--"
		}

		lectureList = append(lectureList, tmp)
	}

	return lectureList
}

func getFavoriteList(r RakutanInfos, pageCount int) []richmenu.FavoritesBodyContents {
	favorites := richmenu.FavoritesJson
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

func getRakutanJudge(r RakutanInfo) RakutanJudge {
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

func intToPtr(i int) *int {
	return &i
}
