package schedule

import (
	"context"
	"event-schedule/internal/model"
)

func (s *Service) UpdateEvent(ctx context.Context, mod *model.UpdateEventInfo) error {
	return s.scheduleRepository.UpdateEvent(ctx, mod)
}
