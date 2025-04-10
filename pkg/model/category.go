package model

type VideoCategory struct {
	BaseModel
	Name   string
	Videos []Video
}
