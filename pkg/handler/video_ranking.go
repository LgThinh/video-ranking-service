package handler

import (
	"github.com/LgThinh/video-ranking-service/pkg/service"
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

	ctx.JSON(http.StatusOK, "")
}

func (h *VideoRankingHandler) GetTopVideoGlobal(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}

func (h *VideoRankingHandler) GetTopVideoPersonalized(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "")
}
