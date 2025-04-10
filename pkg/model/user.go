package model

import (
	"github.com/google/uuid"
)

type User struct {
	BaseModel
	Name        string
	Email       string
	Preferences []UserPreference
}

type UserPreference struct {
	BaseModel
	UserID     uuid.UUID
	User       User
	CategoryID uuid.UUID
	Category   VideoCategory
	Priority   int
}
