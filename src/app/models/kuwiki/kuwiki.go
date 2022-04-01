package kuwiki

// This file was generated from JSON Schema using quicktype, do not modify it directly.

import "encoding/json"

func UnmarshalKUWiki(data []byte) (KUWiki, error) {
	var r KUWiki
	err := json.Unmarshal(data, &r)
	return r, err
}

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
