package schedule

import (
	"context"
	"event-schedule/internal/model"

	"github.com/gofrs/uuid"
)

func (s *Service) AddEvent(ctx context.Context, mod *model.Event) (uuid.UUID, error) {
	return s.scheduleRepository.AddEvent(ctx, mod)
}
