package sync_db

import (
	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"time"
)

func CronSync(db *gorm.DB, client *redis.Client) error {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Hour().Do(SyncToPostgres, db, client)
	if err != nil {
		log.Printf("fail to start cron hourly: %s", err)
		return err
	}
	log.Println("Start cron hourly!")
	s.StartAsync()
	return nil
}

func SyncToPostgres(db *gorm.DB, client *redis.Client) error {
	return nil
}
