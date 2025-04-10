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

type UpdateScoreVideo struct {
	Views     int `json:"views"`
	Likes     int `json:"likes"`
	Comments  int `json:"comments"`
	Shares    int `json:"shares"`
	WatchTime int `json:"watch_time"`
}
