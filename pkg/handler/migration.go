package handler

import (
	"github.com/LgThinh/video-ranking-service/pkg/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"net/http"
)

type MigrationHandler struct {
	db *gorm.DB
}

func NewMigrationHandler(db *gorm.DB) *MigrationHandler {
	return &MigrationHandler{db: db}
}

func (h *MigrationHandler) MigratePublic(ctx *gin.Context) {
	_ = h.db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	_ = h.db.Exec(`SET client_encoding = 'UTF8'`)

	sqlCommands := []string{
		`CREATE SCHEMA IF NOT EXISTS public`,
	}

	for _, sql := range sqlCommands {
		if err := h.db.Exec(sql).Error; err != nil {
			log.Printf(err.Error())
			h.db.Rollback()
			ctx.JSON(http.StatusInternalServerError, err)
		}
	}

	h.db.Config.NamingStrategy = schema.NamingStrategy{
		TablePrefix: "public.",
	}

	models := []interface{}{
		&model.Entity{},
		&model.EntityPreference{},
		&model.Video{},
		&model.VideoCategory{},
		&model.Interaction{},
	}
	if err := h.db.AutoMigrate(models...); err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, "Migrate public schema success")
}
