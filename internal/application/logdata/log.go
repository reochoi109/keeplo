package logdata

import (
	"context"
	"keeplo/internal/adapter/rest/dto"
	"keeplo/internal/domain/logdata"
	"time"
)

type Service struct {
	repo logdata.Repository
}

func NewService(repo logdata.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetResponseTimeSeries(ctx context.Context, monitorID string, from, to time.Time, interval time.Duration, statusCode *int, isSuccess *bool) ([]dto.TimeSeriesPointResponse, error) {
	_, err := s.repo.GetHealthLogs(ctx, monitorID, from, to, statusCode, isSuccess)
	if err != nil {
		return nil, err
	}

	return []dto.TimeSeriesPointResponse{ /* ... */ }, nil
}
