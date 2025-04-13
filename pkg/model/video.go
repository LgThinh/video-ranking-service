package model

import "github.com/google/uuid"

type Video struct {
	BaseModel
	Title      string    `json:"title" gorm:"type:varchar(255);not null"`
	Length     int       `json:"length" gorm:"type:int;default:0"`
	Views      int       `json:"views" gorm:"type:int;default:0"`
	Likes      int       `json:"likes" gorm:"type:int;default:0"`
	Comments   int       `json:"comments" gorm:"type:int;default:0"`
	Shares     int       `json:"shares" gorm:"type:int;default:0"`
	WatchTime  int       `json:"watch_time" gorm:"type:int;default:0"`
	Score      float64   `json:"score" gorm:"type:decimal(15,2);default:0.0"`
	CategoryID uuid.UUID `json:"category_id" gorm:"type:uuid;not null"`
}

func (m *Video) TableName() string {
	return "video"
}

type UpdateScoreVideo struct {
	Views     *int `json:"views"`
	Likes     *int `json:"likes"`
	Comments  *int `json:"comments"`
	Shares    *int `json:"shares"`
	WatchTime *int `json:"watch_time"`
}

type ScoredVideo struct {
	Video      Video   `json:"video"`
	FinalScore float64 `json:"final_score"`
}
