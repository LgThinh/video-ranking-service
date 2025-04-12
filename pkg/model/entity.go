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

type EntityPreference struct {
	BaseModel
	Priority   int       `gorm:"type:int;not null"`
	EntityID   uuid.UUID `gorm:"type:uuid;not null"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null"`
}
