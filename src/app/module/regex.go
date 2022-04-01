package module

import (
	"fmt"
	"regexp"
	"strconv"
)

type PostbackParam struct {
	Type        PostbackType
	ID          int
	LectureName string
}

type PostbackEntry struct {
	Uid   string
	Param PostbackParam
}

type PostbackType string

const (
	Fav = "fav"
	Del = "del"
)

func (pt PostbackType) Valid() error {
	switch pt {
	case Fav, Del:
		return nil
	default:
		return fmt.Errorf("PostbackType: invalid type %s", pt)
	}
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

func IsStudentAddress(address string) bool {
	assigned := regexp.MustCompile(`[A-Za-z0-9._+]+@st\.kyoto-u\.ac\.jp$`)
	return assigned.MatchString(address)
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
		err := PostbackType(matches[0][2]).Valid()
		if err == nil {
			id, _ := strconv.Atoi(matches[1][2])
			params = PostbackParam{Type: PostbackType(matches[0][2]), ID: id, LectureName: matches[2][2]}
			return true, params
		}
	}
	return false, params
}
