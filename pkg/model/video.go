package model

import "github.com/google/uuid"

type Video struct {
	BaseModel
	Title      string    `json:"title"`
	Length     int       `json:"length"`
	Views      int       `json:"views"`
	Likes      int       `json:"likes"`
	Comments   int       `json:"comments"`
	Shares     int       `json:"shares"`
	WatchTime  int       `json:"watch_time"`
	Score      float64   `json:"score"`
	CategoryID uuid.UUID `json:"category_id"`
}

type UpdateScoreVideo struct {
	Views     int `json:"views"`
	Likes     int `json:"likes"`
	Comments  int `json:"comments"`
	Shares    int `json:"shares"`
	WatchTime int `json:"watch_time"`
}
