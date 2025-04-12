package repo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
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
	UpdateGlobalRanking(ctx context.Context, videoID uuid.UUID, score float64)
	IncreaseEntityPriority(ctx context.Context, entityID uuid.UUID, categoryID string)
	GetGlobalRanking(ctx context.Context, limit int64) []redis.Z
	GetPriority(ctx context.Context, entityID, categoryID string) int
	GetCategoryIDByVideoID(ctx context.Context, videoID string) (string, error)
}

func (r *VideoRankingRepo) UpdateGlobalRanking(ctx context.Context, videoID uuid.UUID, score float64) {
	r.Redis.ZAdd(ctx, "ranking:global", redis.Z{
		Score:  score,
		Member: videoID.String(),
	})
}

func (r *VideoRankingRepo) IncreaseEntityPriority(ctx context.Context, entityID uuid.UUID, categoryID string) {
	priorityKey := fmt.Sprintf("priority:%s:%s", entityID.String(), categoryID)
	r.Redis.IncrBy(ctx, priorityKey, 1)
}

func (r *VideoRankingRepo) GetGlobalRanking(ctx context.Context, limit int64) []redis.Z {
	res, _ := r.Redis.ZRevRangeWithScores(ctx, "ranking:global", 0, limit-1).Result()
	return res
}

func (r *VideoRankingRepo) GetPriority(ctx context.Context, entityID, categoryID string) int {
	priorityKey := fmt.Sprintf("priority:%s:%s", entityID, categoryID)
	priorityStr, _ := r.Redis.Get(ctx, priorityKey).Result()
	priority, _ := strconv.Atoi(priorityStr)
	return priority
}

func (r *VideoRankingRepo) GetCategoryIDByVideoID(ctx context.Context, videoID string) (string, error) {
	return r.Redis.HGet(ctx, "video:"+videoID, "category_id").Result()
}
