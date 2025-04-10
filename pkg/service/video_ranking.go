package service

import "github.com/LgThinh/video-ranking-service/pkg/repo"

type VideoRankingService struct {
	videoRankingRepo repo.VideoRankingRepoInterface
}

func NewVideoRankingService(videoRankingRepo repo.VideoRankingRepoInterface) *VideoRankingService {
	return &VideoRankingService{
		videoRankingRepo: videoRankingRepo,
	}
}

type VideoRankingServiceInterface interface{}
