package model

type VideoCategory struct {
	BaseModel
	Name string `json:"name" gorm:"type:varchar(255);not null"`
}
