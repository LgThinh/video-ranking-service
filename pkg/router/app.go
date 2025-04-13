package router

import (
	"context"
	"fmt"
	"github.com/LgThinh/video-ranking-service/conf"
	ginSwaggerDocs "github.com/LgThinh/video-ranking-service/docs"
	cron "github.com/LgThinh/video-ranking-service/pkg/cron/sync_db"
	handlers "github.com/LgThinh/video-ranking-service/pkg/handler"
	"github.com/LgThinh/video-ranking-service/pkg/middlewares"
	"github.com/LgThinh/video-ranking-service/pkg/repo"
	"github.com/LgThinh/video-ranking-service/pkg/service"
	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"net/http"
	"strings"
	"time"
)

func NewRouter() {
	router := gin.Default()

	configCors := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "x-entity-id"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	db := initPostgres()
	redisClient := initRedis()

	router.Use(cors.New(configCors))
	router.Use(limit.MaxAllowed(200))

	err := cron.CronSync(db, redisClient)
	if err != nil {
		log.Println(err)
	}

	ApplicationV1Router(router, db, redisClient)
	startServer(router)
}

func ApplicationV1Router(router *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	// Router
	routerV1 := router.Group("/api/v1")

	// Init repo
	videoRankingRepo := repo.NewVideoRankingRepo(db, redisClient)

	// Init service
	videoRankingService := service.NewVideoRankingService(videoRankingRepo)

	// Init handler
	migrateHandler := handlers.NewMigrationHandler(db)
	videoRankingHandler := handlers.NewVideoRankingHandler(videoRankingService)

	// Migrate api
	routerV1.POST("/internal/migrate", middlewares.AuthJWTMiddleware(), migrateHandler.MigratePublic)

	// Video Ranking Apis
	routerV1.PUT("/score/update", videoRankingHandler.UpdateVideoScore)
	routerV1.PUT("/entity-preference/update", videoRankingHandler.UpdateEntityPreference)
	routerV1.GET("/video-global/list", videoRankingHandler.GetTopVideoGlobal)
	routerV1.GET("/video-personalized/list", videoRankingHandler.GetTopVideoPersonalized)

	// Swagger
	ginSwaggerDocs.SwaggerInfo.Host = conf.GetConfig().SwaggerHost
	ginSwaggerDocs.SwaggerInfo.Title = conf.GetConfig().AppName
	ginSwaggerDocs.SwaggerInfo.BasePath = "/api/v1"
	ginSwaggerDocs.SwaggerInfo.Version = "v1"
	ginSwaggerDocs.SwaggerInfo.Schemes = []string{"http", "https"}

	routerV1.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.PersistAuthorization(true),
	))

}

func startServer(router http.Handler) {
	s := &http.Server{
		Addr:           ":" + conf.GetConfig().Port,
		Handler:        router,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Server running on port %s", conf.GetConfig().Port)
	if err := s.ListenAndServe(); err != nil {
		_ = fmt.Errorf("fatal error description: %s", strings.ToLower(err.Error()))
		panic(err)
	}
}

func initPostgres() *gorm.DB {
	dsn := postgres.Open(fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable connect_timeout=5",
		conf.GetConfig().DBHost,
		conf.GetConfig().DBPort,
		conf.GetConfig().DBUser,
		conf.GetConfig().DBName,
		conf.GetConfig().DBPass,
	))
	db, err := gorm.Open(dsn, &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{
			SingularTable: true,
			//TablePrefix:   conf.GetConfig().DBSchema + ".",
		},
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		log.Fatalf("error opening connection to database: %v", err)
	}

	conn, err := db.DB()
	if err != nil {
		log.Fatalf("error initializing database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err = conn.PingContext(ctx); err != nil {
		log.Fatalf("error opening connection to database: %v", err)
	}

	log.Printf("Postgres connected!")
	return db
}

func initRedis() *redis.Client {
	options := &redis.Options{
		Addr:     conf.GetConfig().RedisAddress,
		Password: conf.GetConfig().RedisPassword,
		DB:       conf.GetConfig().RedisDB,
	}

	client := redis.NewClient(options)

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Redis connected!")
	return client
}
