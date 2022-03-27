package models

type RakutanInfo struct {
	ID          int             `bson:"id"`
	FacultyName string          `bson:"faculty_name"`
	LectureName string          `bson:"lecture_name"`
	Detail      []RakutanDetail `bson:"detail"`
	URL         string          `bson:"url,omitempty"`
}

type RakutanDetail struct {
	Accepted int `bson:"accepted,omitempty"`
	Total    int `bson:"total,omitempty"`
}
