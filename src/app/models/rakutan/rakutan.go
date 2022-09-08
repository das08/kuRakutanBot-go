package models

import "github.com/lib/pq"

type RakutanInfo struct {
	ID          int            `bson:"id"`
	FacultyName string         `bson:"faculty_name"`
	LectureName string         `bson:"lecture_name"`
	Detail      RakutanDetails `bson:"detail"`
	OmikujiType string         `bson:"omikuji,omitempty"`
	URL         string         `bson:"url,omitempty"`
	IsVerified  bool
	IsFavorite  bool
	KUWikiErr   string
}

type RakutanInfo2 struct {
	ID          int           `db:"id"`
	FacultyName string        `db:"faculty_name"`
	LectureName string        `db:"lecture_name"`
	Register    pq.Int32Array `db:"register"`
	Passed      pq.Int32Array `db:"passed"`
}

type RakutanDetail struct {
	Year     int `bson:"year,omitempty"`
	Accepted int `bson:"accepted,omitempty"`
	Total    int `bson:"total,omitempty"`
}

type RakutanDetails []RakutanDetail

func (rds RakutanDetails) GetLatestDetail() (int, int) {
	accept, total := 0, 0
	for _, rd := range rds {
		if rd.Total != 0 {
			accept = rd.Accepted
			total = rd.Total
			break
		}
	}
	return accept, total
}

type UserData struct {
	Uid          string    `bson:"uid"`
	Count        UserCount `bson:"count"`
	RegisterTime int       `bson:"register_time"`
	Verified     bool      `bson:"verified"`
}

type UserCount struct {
	Message        int `bson:"message,omitempty"`
	RakutanOmikuji int `bson:"rakutan_omikuji,omitempty"`
	OnitanOmikuji  int `bson:"onitan_omikuji,omitempty"`
}
