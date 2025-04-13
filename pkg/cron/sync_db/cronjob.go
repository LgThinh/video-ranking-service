package sync_db

import (
	"context"
	"fmt"
	"github.com/LgThinh/video-ranking-service/pkg/model"
	"github.com/LgThinh/video-ranking-service/pkg/utils"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"time"
)

func CronSync(db *gorm.DB, client *redis.Client) error {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Minute().Do(SyncToPostgres, db, client)
	if err != nil {
		log.Printf("fail to start cron hourly: %s", err)
	}
	log.Println("Start cron hourly!")
	s.StartAsync()
	return nil
}

func SyncToPostgres(db *gorm.DB, client *redis.Client) error {
	ctx := context.Background()
	topVideos, err := client.ZRevRangeWithScores(ctx, "video_score_ranking", 0, -1).Result()
	if err != nil {
		return err
	}

	for _, z := range topVideos {
		videoIDStr := z.Member.(string)
		videoID, err := uuid.Parse(videoIDStr)
		if err != nil {
			log.Println("error parse video ID", videoIDStr, ":", err)
			continue
		}

		statsKey := fmt.Sprintf("video_stats:%s", videoIDStr)
		stats, err := client.HGetAll(ctx, statsKey).Result()
		if err != nil {
			log.Println("error get video stats for", videoIDStr, ":", err)
			continue
		}

		views := utils.ParseInt(stats["views"])
		likes := utils.ParseInt(stats["likes"])
		comments := utils.ParseInt(stats["comments"])
		shares := utils.ParseInt(stats["shares"])
		watchTime := utils.ParseInt(stats["watch_time"])

		err = UpdateScorePostgres(db, videoID, z.Score, views, likes, comments, shares, watchTime)
		if err != nil {
			log.Println("Update score failed for", videoIDStr, ":", err)
			continue
		}
	}
	return nil
}

func UpdateScorePostgres(db *gorm.DB, videoID uuid.UUID, score float64, views, likes, comments, shares, watchTime int) error {
	return db.Model(&model.Video{}).Where("id = ?", videoID).Updates(map[string]interface{}{
		"score":      score,
		"views":      views,
		"likes":      likes,
		"comments":   comments,
		"shares":     shares,
		"watch_time": watchTime,
	}).Error
}
