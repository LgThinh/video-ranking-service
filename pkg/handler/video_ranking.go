package handler

import (
	"github.com/LgThinh/video-ranking-service/pkg/model"
	"github.com/LgThinh/video-ranking-service/pkg/service"
	"github.com/LgThinh/video-ranking-service/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type VideoRankingHandler struct {
	videoRankingService service.VideoRankingServiceInterface
}

func NewVideoRankingHandler(videoRankingService service.VideoRankingServiceInterface) *VideoRankingHandler {
	return &VideoRankingHandler{videoRankingService: videoRankingService}
}

func (h *VideoRankingHandler) UpdateVideoScore(ctx *gin.Context) {
	var req model.UpdateScoreVideo
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	videoIDStr := ctx.GetHeader("x-video-id")
	videoID, err := utils.ParseIDtoUUID(videoIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	entityIDStr := ctx.GetHeader("x-entity-id")
	entityID, err := utils.ParseIDtoUUID(entityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	err = h.videoRankingService.UpdateVideoScore(ctx, *videoID, *entityID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	ctx.JSON(http.StatusOK, "Update video ranking success")
}

func (h *VideoRankingHandler) GetTopVideoGlobal(ctx *gin.Context) {
	top := h.videoRankingService.GetGlobalRanking(ctx)

	ctx.JSON(http.StatusOK, top)
}

func (h *VideoRankingHandler) GetTopVideoPersonalized(ctx *gin.Context) {
	entityID := ctx.Param("entity_id")
	top := h.videoRankingService.GetEntityRanking(ctx, entityID)

	ctx.JSON(http.StatusOK, top)
}
