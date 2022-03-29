package module

import (
	"regexp"
	"strconv"
)

type PostbackParam struct {
	Type        string
	ID          int
	LectureName string
}

type PostbackEntry struct {
	Uid   string
	Param PostbackParam
}

func IsLectureID(messageText string) (bool, int) {
	assigned := regexp.MustCompile("^#([1-9][0-9]{4})$")
	group := assigned.FindSubmatch([]byte(messageText))
	if len(group) == 2 {
		id, _ := strconv.Atoi(string(group[1]))
		return true, id
	}
	return false, 0
}

func ParsePBParam(messageText string) (bool, PostbackParam) {
	expectedKey := [3]string{"type", "id", "lecname"}
	assigned := regexp.MustCompile("([^=&]+)=([^&]*)")
	matches := assigned.FindAllStringSubmatch(messageText, -1)
	params := PostbackParam{}
	if len(matches) == 3 {
		for i, p := range matches {
			if p[1] != expectedKey[i] {
				return false, params
			}
		}
		id, _ := strconv.Atoi(matches[1][2])
		params = PostbackParam{Type: matches[0][2], ID: id, LectureName: matches[2][2]}
		return true, params
	}
	return false, params
}
