package model

import "github.com/google/uuid"

type Interaction struct {
	BaseModel
	EntityID  uuid.UUID `json:"entity_id" gorm:"type:uuid;not null"`
	VideoID   uuid.UUID `json:"video_id" gorm:"type:uuid;not null"`
	View      bool      `json:"view" gorm:"default:false"`
	Like      bool      `json:"like" gorm:"default:false"`
	Commented bool      `json:"commented" gorm:"default:false"`
	Shared    bool      `json:"shared" gorm:"default:false"`
	WatchTime int       `json:"watch_time" gorm:"type:int;default:0"`
}

func (m *Interaction) TableName() string {
	return "interaction"
}
