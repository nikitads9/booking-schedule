package schedule

import (
	"context"
	"event-schedule/internal/model"
)

func (s *Service) GetEvent(ctx context.Context, eventID string) (*model.Event, error) {
	return s.scheduleRepository.GetEvent(ctx, eventID)
}
