package model

import (
	"github.com/google/uuid"
)

type Entity struct {
	BaseModel
	Name        string `json:"name" gorm:"type:varchar(255);not null"`
	UserName    string `json:"user_name" gorm:"type:varchar(255);unique;not null"`
	Password    string `json:"password" gorm:"type:varchar(255);not null"`
	Email       string `json:"email" gorm:"type:varchar(255);unique;not null"`
	PhoneNumber string `json:"phone_number" gorm:"type:varchar(255);unique;not null"`
	EntityType  string `json:"entity_type" gorm:"type:varchar(255);not null"`
}

func (m *Entity) TableName() string {
	return "entity"
}

type EntityPreference struct {
	BaseModel
	Priority   float64   `json:"priority" gorm:"type:decimal(15,2);default:0.0"`
	EntityID   uuid.UUID `json:"entity_id" gorm:"type:uuid;not null"`
	CategoryID uuid.UUID `json:"category_id" gorm:"type:uuid;not null"`
}

func (m *EntityPreference) TableName() string {
	return "entity_preference"
}
