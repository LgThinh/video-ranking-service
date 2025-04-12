package repo

import (
	"context"
	"github.com/LgThinh/video-ranking-service/pkg/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type VideoRankingRepo struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func (r *VideoRankingRepo) BeginTransaction() *gorm.DB {
	return r.DB.Begin()
}

func (r *VideoRankingRepo) DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout)
	return r.DB.WithContext(ctx), cancel
}

func NewVideoRankingRepo(videoRankingRepo *gorm.DB, redisClient *redis.Client) VideoRankingRepoInterface {
	return &VideoRankingRepo{
		DB:    videoRankingRepo,
		Redis: redisClient,
	}
}

type VideoRankingRepoInterface interface {
	BeginTransaction() *gorm.DB
	DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc)
	UpdateVideoScore(tx *gorm.DB, videoID uuid.UUID, score float64) error
	GetVideoByID(tx *gorm.DB, videoID uuid.UUID) (*model.Video, error)
	GetEntityPreference(tx *gorm.DB, entityID, categoryID uuid.UUID) (*model.EntityPreference, error)
	UpdateEntityPreference(tx *gorm.DB, entityID, categoryID uuid.UUID, priority float64) error
	GetTopVideoGlobal(tx *gorm.DB) (*[]model.Video, error)
}

func (r *VideoRankingRepo) UpdateVideoScore(tx *gorm.DB, videoID uuid.UUID, score float64) error {
	scoreUpdate := map[string]interface{}{
		"score":      score,
		"updated_at": time.Now(),
	}

	err := tx.Model(&model.Video{}).Where("id = ?", videoID).Updates(&scoreUpdate).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *VideoRankingRepo) GetVideoByID(tx *gorm.DB, videoID uuid.UUID) (*model.Video, error) {
	var video model.Video

	err := tx.Model(&model.Video{}).Where("id = ?", videoID).First(&video).Error
	if err != nil {
		return nil, err
	}

	return &video, nil
}

func (r *VideoRankingRepo) GetEntityPreference(tx *gorm.DB, entityID, categoryID uuid.UUID) (*model.EntityPreference, error) {
	var entityPreference model.EntityPreference

	err := tx.Model(&model.EntityPreference{}).Where("entity_id = ? AND category_id = ?", entityID, categoryID).First(&entityPreference).Error
	if err != nil {
		return nil, err
	}

	return &entityPreference, nil
}

func (r *VideoRankingRepo) UpdateEntityPreference(tx *gorm.DB, entityID, categoryID uuid.UUID, priority float64) error {
	var p model.EntityPreference

	result := tx.First(&p, "entity_id = ? AND category_id = ?", entityID, categoryID)
	if result.Error != nil {
		tx.Create(&model.EntityPreference{
			EntityID:   entityID,
			CategoryID: categoryID,
			Priority:   priority,
		})
	} else {
		tx.Model(&model.EntityPreference{}).
			Where("entity_id = ? AND category_id = ?", entityID, categoryID).
			Update("priority", gorm.Expr("priority + ?", priority))
	}

	return nil
}

func (r *VideoRankingRepo) GetTopVideoGlobal(tx *gorm.DB) (*[]model.Video, error) {
	var videos []model.Video

	err := tx.Model(&model.Video{}).Order("score desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}

	return &videos, nil
}
