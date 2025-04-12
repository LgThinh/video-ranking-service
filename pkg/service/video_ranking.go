package service

import (
	"context"
	"github.com/LgThinh/video-ranking-service/pkg/model"
	"github.com/LgThinh/video-ranking-service/pkg/repo"
	"github.com/LgThinh/video-ranking-service/pkg/utils"
	"github.com/google/uuid"
	"sort"
)

type VideoRankingService struct {
	videoRankingRepo repo.VideoRankingRepoInterface
}

func NewVideoRankingService(videoRankingRepo repo.VideoRankingRepoInterface) *VideoRankingService {
	return &VideoRankingService{
		videoRankingRepo: videoRankingRepo,
	}
}

type VideoRankingServiceInterface interface {
	UpdateVideoScore(ctx context.Context, videoID, entityID uuid.UUID, req model.UpdateScoreVideo) error
	GetGlobalRanking(ctx context.Context) []string
	GetEntityRanking(ctx context.Context, entityID string) []string
}

func (s *VideoRankingService) UpdateVideoScore(ctx context.Context, videoID, entityID uuid.UUID, req model.UpdateScoreVideo) error {
	score := utils.CalculateScore(req)
	s.videoRankingRepo.UpdateGlobalRanking(ctx, videoID, score)
	return nil
}

func (s *VideoRankingService) GetGlobalRanking(ctx context.Context) []string {
	videos := s.videoRankingRepo.GetGlobalRanking(ctx, 10)
	ids := []string{}
	for _, z := range videos {
		ids = append(ids, z.Member.(string))
	}
	return ids
}

func (s *VideoRankingService) GetEntityRanking(ctx context.Context, entityID string) []string {
	videos := s.videoRankingRepo.GetGlobalRanking(ctx, 100)
	result := make([]struct {
		ID   string
		Real float64
	}, 0)

	for _, z := range videos {
		videoID := z.Member.(string)

		categoryID, err := s.videoRankingRepo.GetCategoryIDByVideoID(ctx, videoID)
		if err != nil {
			continue
		}
		priority := s.videoRankingRepo.GetPriority(ctx, entityID, categoryID)
		finalScore := z.Score * (1 + float64(priority)*0.01)
		result = append(result, struct {
			ID   string
			Real float64
		}{videoID, finalScore})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Real > result[j].Real
	})

	top := []string{}
	for _, v := range result[:10] {
		top = append(top, v.ID)
	}
	return top
}
