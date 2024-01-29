package schedule

import (
	"context"
	"event-schedule/internal/model"

	"github.com/gofrs/uuid"
)

func (s *Service) GetEvent(ctx context.Context, eventID uuid.UUID) (*model.EventInfo, error) {
	return s.scheduleRepository.GetEvent(ctx, eventID)
}
