package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/LgThinh/video-ranking-service/pkg/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
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
	UpdateScoreInRedis(ctx context.Context, videoID uuid.UUID, score float64) error
	UpdateStatsInRedis(ctx context.Context, videoID uuid.UUID, req model.UpdateScoreVideo) error
	GetVideoByID(tx *gorm.DB, videoID uuid.UUID) (*model.Video, error)
	GetEntityPreference(tx *gorm.DB, entityID, categoryID uuid.UUID) (*model.EntityPreference, error)
	UpdateEntityPreference(tx *gorm.DB, entityID, categoryID uuid.UUID, priority float64) error
	GetTopVideoGlobal(tx *gorm.DB) (*[]model.Video, error)
	UpsertInteraction(tx *gorm.DB, entityID, videoID uuid.UUID, req *model.UpdateEntityPreference) error
}

func (r *VideoRankingRepo) UpdateScoreInRedis(ctx context.Context, videoID uuid.UUID, score float64) error {
	return r.Redis.ZIncrBy(ctx, "video_score_ranking", score, videoID.String()).Err()
}

func (r *VideoRankingRepo) UpdateStatsInRedis(ctx context.Context, videoID uuid.UUID, req model.UpdateScoreVideo) error {
	key := fmt.Sprintf("video_stats:%s", videoID.String())

	pipe := r.Redis.TxPipeline()

	if req.Views != nil && *req.Views > 0 {
		pipe.HIncrBy(ctx, key, "views", int64(*req.Views))
	}
	if req.Likes != nil && *req.Likes > 0 {
		pipe.HIncrBy(ctx, key, "likes", int64(*req.Likes))
	}
	if req.Comments != nil && *req.Comments > 0 {
		pipe.HIncrBy(ctx, key, "comments", int64(*req.Comments))
	}
	if req.Shares != nil && *req.Shares > 0 {
		pipe.HIncrBy(ctx, key, "shares", int64(*req.Shares))
	}
	if req.WatchTime != nil && *req.WatchTime > 0 {
		pipe.HIncrBy(ctx, key, "watch_time", int64(*req.WatchTime))
	}

	_, err := pipe.Exec(ctx)
	return err
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

func (r *VideoRankingRepo) UpsertInteraction(tx *gorm.DB, entityID, videoID uuid.UUID, req *model.UpdateEntityPreference) error {
	var interaction model.Interaction

	err := tx.Where("entity_id = ? AND video_id = ?", entityID, videoID).First(&interaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newInteraction := model.Interaction{
				EntityID:  entityID,
				VideoID:   videoID,
				View:      req.Views != nil && *req.Views,
				Like:      req.Likes != nil && *req.Likes,
				Commented: req.Comments != nil && *req.Comments,
				Shared:    req.Shares != nil && *req.Shares,
				WatchTime: 0,
			}
			if req.WatchTime != nil && *req.WatchTime > 0 {
				newInteraction.WatchTime = *req.WatchTime
			}
			return tx.Create(&newInteraction).Error
		}
		return err
	}

	if req.Views != nil {
		interaction.View = *req.Views
	}
	if req.Likes != nil {
		interaction.Like = *req.Likes
	}
	if req.Comments != nil {
		interaction.Commented = *req.Comments
	}
	if req.Shares != nil {
		interaction.Shared = *req.Shares
	}
	if req.WatchTime != nil {
		interaction.WatchTime += *req.WatchTime
	}

	return tx.Save(&interaction).Error
}
