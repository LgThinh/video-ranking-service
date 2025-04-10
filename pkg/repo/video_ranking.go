package repo

import (
	"context"
	"github.com/go-redis/redis/v8"
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

type VideoRankingRepoInterface interface{}
