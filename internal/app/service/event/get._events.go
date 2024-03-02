package event

import (
	"context"
	"event-schedule/internal/app/model"
	"time"
)

func (s *Service) GetEvents(ctx context.Context, startDate time.Time, endDate time.Time, id int64) ([]*model.EventInfo, error) {
	return s.eventRepository.GetEvents(ctx, startDate, endDate, id)
}
