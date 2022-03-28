package module

import (
	"regexp"
	"strconv"
)

func IsLectureNumber(messageText string) (bool, int) {
	assined := regexp.MustCompile("^#([1-9][0-9]{5})$")
	group := assined.FindSubmatch([]byte(messageText))
	if len(group) == 2 {
		id, _ := strconv.Atoi(string(group[1]))
		return true, id
	}
	return false, 0
}
