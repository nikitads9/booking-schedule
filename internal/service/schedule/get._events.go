package schedule

import (
	"context"
	"event-schedule/internal/model"
)

func (s *Service) GetEvents(ctx context.Context, mod *model.GetEventsInfo) ([]*model.EventInfo, error) {
	return s.scheduleRepository.GetEvents(ctx, mod)
}
