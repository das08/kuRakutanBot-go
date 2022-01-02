package models

type RakutanInfo struct {
	ID          int           `bson:"id"`
	FacultyName string        `bson:"facultyname"`
	LectureName string        `bson:"lecturename"`
	Groups      string        `bson:"groups,omitempty"`
	Credits     string        `bson:"credits,omitempty"`
	Detail      RakutanDetail `bson:"detail,omitempty"`
	URL         string        `bson:"url,omitempty"`
}

type RakutanDetail struct {
	Accepted int `bson:"accepted,omitempty"`
	Total    int `bson:"total,omitempty"`
}
