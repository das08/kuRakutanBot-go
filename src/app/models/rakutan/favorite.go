package models

type Favorite struct {
	Uid         string `bson:"uid,omitempty"`
	ID          int    `bson:"id,omitempty"`
	LectureName string `bson:"lecture_name,omitempty"`
}
