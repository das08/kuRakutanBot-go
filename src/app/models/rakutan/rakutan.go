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

func (rds RakutanDetails) AcceptedList() []int {
	var list []int
	for _, rd := range rds {
		list = append(list, rd.Accepted)
	}
	return list
}

func (rds RakutanDetails) TotalList() []int {
	var list []int
	for _, rd := range rds {
		list = append(list, rd.Total)
	}
	return list
}
