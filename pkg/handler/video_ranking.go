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

// UpdateVideoScore godoc
// @Summary     Update Video Score
// @Tags        Video Ranking
// @Accept      json
// @Produce     json
// @Param       x-video-id header string true "VideoID"
// @Param       update_score body model.UpdateScoreVideo true "UpdateScoreVideo"
// @Success     200 {string} string "Update video ranking success"
// @Failure     400 {object} map[string]interface{}
// @Router		/score/update [put]
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
		return
	}

	err = h.videoRankingService.UpdateVideoScore(ctx, *videoID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, "Update video ranking success")
}

// UpdateEntityPreference godoc
// @Summary     Update Entity Preference
// @Tags        Video Ranking
// @Accept      json
// @Produce     json
// @Param       x-video-id header string true "VideoID"
// @Param       x-entity-id header string	true "EntityID"
// @Param       update_score body model.UpdateEntityPreference true "UpdateEntityPreference"
// @Success     200 {string} string "Update entity preference success"
// @Failure     400 {object} map[string]interface{}
// @Router		/entity-preference/update [put]
func (h *VideoRankingHandler) UpdateEntityPreference(ctx *gin.Context) {
	var req model.UpdateEntityPreference
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	videoIDStr := ctx.GetHeader("x-video-id")
	videoID, err := utils.ParseIDtoUUID(videoIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	entityIDStr := ctx.GetHeader("x-entity-id")
	entityID, err := utils.ParseIDtoUUID(entityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = h.videoRankingService.UpdateEntityPreference(ctx, *videoID, *entityID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, "Update entity preference success")
}

// GetTopVideoGlobal godoc
// @Summary     Get Top Global Videos
// @Tags        Video Ranking
// @Accept      json
// @Produce     json
// @Success     200 {array} model.Video
// @Failure     400 {object} map[string]interface{}
// @Router      /video-global/list [get]
func (h *VideoRankingHandler) GetTopVideoGlobal(ctx *gin.Context) {
	topVideos, err := h.videoRankingService.GetGlobalRanking(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, topVideos)
}

// GetTopVideoPersonalized godoc
// @Summary     Get Personalized Top Videos
// @Tags        Video Ranking
// @Accept      json
// @Produce     json
// @Param       x-entity-id header string	true "EntityID"
// @Success     200 {array} model.Video
// @Failure     400 {object} map[string]interface{}
// @Router      /video-personalized/list [get]
func (h *VideoRankingHandler) GetTopVideoPersonalized(ctx *gin.Context) {
	entityIDStr := ctx.GetHeader("x-entity-id")
	entityID, err := utils.ParseIDtoUUID(entityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	topVideos, err := h.videoRankingService.GetEntityRanking(ctx, *entityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, topVideos)
}
