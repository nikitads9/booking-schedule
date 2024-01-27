package schedule

import (
	"context"
	"event-schedule/internal/model"
)

func (s *Service) GetEvents(ctx context.Context, userID int64, period string) ([]*model.Event, error) {
	return s.scheduleRepository.GetEvents(ctx, userID)
}
