package utils

import (
	"fmt"
	"github.com/LgThinh/video-ranking-service/pkg/model"
	"github.com/google/uuid"
	"log"
	"strconv"
)

func ParseIDtoUUID(id interface{}) (*uuid.UUID, error) {
	if id == nil {
		return nil, fmt.Errorf("id cannot be nil")
	}

	newIDStr, ok := id.(string)
	if !ok {
		return nil, fmt.Errorf("id is not a string")
	}
	newUUID, err := uuid.Parse(newIDStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &newUUID, nil
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func BoolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0
}

func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func CalculateScore(interaction model.UpdateScoreVideo) float64 {
	return float64(*interaction.Views)*1.0 +
		float64(*interaction.Likes)*1.5 +
		float64(*interaction.Comments)*2.0 +
		float64(*interaction.Shares)*3.0 +
		float64(*interaction.WatchTime)*0.1
}

func CalculatePriority(interaction model.UpdateEntityPreference) float64 {
	var (
		views       float64
		likes       float64
		comments    float64
		shares      float64
		watchTime   float64
		newPriority float64
	)
	if interaction.Views != nil {
		views = BoolToFloat64(*interaction.Views)
	}
	if interaction.Likes != nil {
		likes = BoolToFloat64(*interaction.Likes)
	}
	if interaction.Comments != nil {
		comments = BoolToFloat64(*interaction.Comments)
	}
	if interaction.Shares != nil {
		shares = BoolToFloat64(*interaction.Shares)
	}
	if interaction.WatchTime != nil {
		watchTime = float64(*interaction.WatchTime)
	}

	newPriority = views + likes*2.0 + comments*3.0 + shares*4.0 + watchTime*0.2

	return newPriority
}
