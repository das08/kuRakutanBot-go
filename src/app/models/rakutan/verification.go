package models

type Verification struct {
	Uid   string `bson:"uid,omitempty"`
	Code  string `bson:"code,omitempty"`
	Email string
}
