package utils

import (
	"fmt"
	"github.com/LgThinh/video-ranking-service/pkg/model"
	"github.com/google/uuid"
	"log"
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

func CalculateScore(interaction model.UpdateScoreVideo) float64 {
	return float64(*interaction.Views)*1.0 +
		float64(*interaction.Likes)*1.5 +
		float64(*interaction.Comments)*2.0 +
		float64(*interaction.Shares)*3.0 +
		float64(*interaction.WatchTime)*0.1
}

func CalculatePriority(interaction model.UpdateScoreVideo) float64 {
	return float64(*interaction.Views)*2.0 +
		float64(*interaction.Likes)*1.0 +
		float64(*interaction.Shares)*3.0
}
