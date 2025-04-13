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
	UpdateVideoScore(ctx context.Context, videoID uuid.UUID, req model.UpdateScoreVideo) error
	UpdateEntityPreference(ctx context.Context, videoID, entityID uuid.UUID, req model.UpdateEntityPreference) error
	GetGlobalRanking(ctx context.Context) (*[]model.Video, error)
	GetEntityRanking(ctx context.Context, entityID uuid.UUID) (*[]model.Video, error)
}

func (s *VideoRankingService) UpdateVideoScore(ctx context.Context, videoID uuid.UUID, req model.UpdateScoreVideo) error {
	score := utils.CalculateScore(req)

	err := s.videoRankingRepo.UpdateScoreInRedis(ctx, videoID, score)
	if err != nil {
		return err
	}

	err = s.videoRankingRepo.UpdateStatsInRedis(ctx, videoID, req)
	if err != nil {
		return err
	}

	return nil
}

func (s *VideoRankingService) UpdateEntityPreference(ctx context.Context, videoID, entityID uuid.UUID, req model.UpdateEntityPreference) error {
	txWithTimeout, cancel := s.videoRankingRepo.DBWithTimeout(ctx)
	defer cancel()

	video, err := s.videoRankingRepo.GetVideoByID(txWithTimeout, videoID)
	if err != nil {
		return nil
	}

	newPriority := utils.CalculatePriority(req)
	err = s.videoRankingRepo.UpdateEntityPreference(txWithTimeout, entityID, video.CategoryID, newPriority)
	if err != nil {
		return nil
	}

	err = s.videoRankingRepo.UpsertInteraction(txWithTimeout, entityID, videoID, &req)
	if err != nil {
		return nil
	}

	return nil
}

func (s *VideoRankingService) GetGlobalRanking(ctx context.Context) (*[]model.Video, error) {
	txWithTimeout, cancel := s.videoRankingRepo.DBWithTimeout(ctx)
	defer cancel()

	videos, err := s.videoRankingRepo.GetTopVideoGlobal(txWithTimeout)
	if err != nil {
		return nil, err
	}

	return videos, nil
}

func (s *VideoRankingService) GetEntityRanking(ctx context.Context, entityID uuid.UUID) (*[]model.Video, error) {
	tx, cancel := s.videoRankingRepo.DBWithTimeout(ctx)
	defer cancel()

	videos, err := s.videoRankingRepo.GetTopVideoGlobal(tx)
	if err != nil {
		return nil, err
	}

	scored := make([]model.ScoredVideo, 0, len(*videos))

	for _, video := range *videos {
		finalScore := video.Score

		// Kiểm tra xem user có preference với category này không
		entityPreference, err := s.videoRankingRepo.GetEntityPreference(tx, entityID, video.CategoryID)
		if err == nil && entityPreference != nil {
			entityPriority := 1.0 + float64(entityPreference.Priority)*0.01
			finalScore = video.Score * entityPriority
		}

		scored = append(scored, model.ScoredVideo{
			Video:      video,
			FinalScore: finalScore,
		})
	}

	// Sắp xếp tất cả video theo finalScore
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].FinalScore > scored[j].FinalScore
	})

	result := make([]model.Video, 0, len(scored))
	for _, v := range scored {
		result = append(result, v.Video)
	}

	return &result, nil
}
