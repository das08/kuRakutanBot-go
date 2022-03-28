package models

type RakutanInfo struct {
	ID          int            `bson:"id"`
	FacultyName string         `bson:"faculty_name"`
	LectureName string         `bson:"lecture_name"`
	Detail      RakutanDetails `bson:"detail"`
	URL         string         `bson:"url,omitempty"`
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
