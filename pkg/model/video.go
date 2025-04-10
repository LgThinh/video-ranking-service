package model

import (
	"github.com/google/uuid"
)

type Video struct {
	BaseModel
	Title      string
	CategoryID uuid.UUID
	Category   VideoCategory
	Views      int
	Likes      int
	Comments   int
	Shares     int
	WatchTime  int
	Score      float64
}
