package module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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

func (r *KUWiki) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type KUWikiStatus struct {
	Success bool
	Result  string
}

func GetKakomonURL(e *Environments, lectureName string) (ExecStatus[KUWikiKakomon], bool) {
	var status ExecStatus[KUWikiKakomon]
	var kakomonURL KUWikiKakomon
	req, err := http.NewRequest("GET", e.KuwikiEndpoint, nil)
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
		status.Err = ErrorMessageKUWikiGetFailed
		return status, false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll err=%s", err.Error())
		status.Err = ErrorMessageKUWikiGetFailed
		return status, false
	}

	response := KUWiki{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		status.Err = ErrorMessageKUWikiGetFailed
		return status, false
	}

	for _, result := range response.Results {
		if result.Name == lectureName {
			for _, exam := range result.ExamSet {
				kakomonURL = KUWikiKakomon(exam.DriveLink)
			}
		}
	}
	log.Println("[KUWiki] Got kakomon URL")
	status.Result = kakomonURL
	return status, true
}
