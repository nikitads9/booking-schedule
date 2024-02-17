package event

import (
	"context"
	"event-schedule/internal/model"
)

func (s *Service) GetVacantDates(ctx context.Context, suiteID int64) ([]*model.Interval, error) {
	return s.eventRepository.GetVacantDates(ctx, suiteID)
}
