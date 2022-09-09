package module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (r *KUWiki) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type KUWiki struct {
	Count    int64       `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []Result    `json:"results"`
}

type Result struct {
	ID              int64        `json:"id"`
	CourseCode      string       `json:"course_code"`
	CourseNumbering string       `json:"course_numbering"`
	Name            string       `json:"name"`
	Field           string       `json:"field"`
	LectureSet      []LectureSet `json:"lecture_set"`
	ExamSet         []ExamSet    `json:"exam_set"`
	ExamCount       int64        `json:"exam_count"`
}

type ExamSet struct {
	CourseCode   string `json:"course_code"`
	Name         string `json:"name"`
	Field        string `json:"field"`
	DriveID      string `json:"drive_id"`
	DriveLink    string `json:"drive_link"`
	DriveLinkTag string `json:"drive_link_tag"`
}

type LectureSet struct {
	Year            int64  `json:"year"`
	GroupCode       string `json:"group_code"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	InstructorSet   []Set  `json:"instructor_set"`
	PeriodSet       []Set  `json:"period_set"`
	Semester        string `json:"semester"`
	Major           string `json:"major"`
	URL             string `json:"url"`
	NumPeriods      int64  `json:"num_periods"`
	CourseCode      string `json:"course_code"`
	CourseNumbering string `json:"course_numbering"`
}

type Set struct {
	Lecture string `json:"lecture"`
	Name    string `json:"name"`
}

type KUWikiStatus struct {
	Success bool
	Result  string
}

func GetKakomonURL(e *Environments, lectureName string) KUWikiStatus {
	kakomonURL := ""
	method := "GET"
	req, err := http.NewRequest(method, e.KuwikiEndpoint, nil)
	if err != nil {
		log.Fatalf("NewRequest err=%s", err.Error())
	}

	q := req.URL.Query()
	q.Add("name", lectureName)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", e.KuwikiAccessToken))

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Client.Do err=%s", err.Error())
		return KUWikiStatus{false, "取得失敗"}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err=%s", err.Error())
		return KUWikiStatus{false, "取得失敗"}
	}

	response := KUWiki{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("json.Unmarshal err=%s", err.Error())
		return KUWikiStatus{false, "取得失敗"}
	}

	for _, result := range response.Results {
		if result.Name == lectureName {
			for _, exam := range result.ExamSet {
				kakomonURL = exam.DriveLink
			}
		}
	}
	log.Println("[KUWiki] Got kakomon URL")
	return KUWikiStatus{true, kakomonURL}
}
