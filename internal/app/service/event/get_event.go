package event

import (
	"context"
	"event-schedule/internal/app/model"

	"github.com/gofrs/uuid"
)

func (s *Service) GetEvent(ctx context.Context, eventID uuid.UUID) (*model.EventInfo, error) {
	return s.eventRepository.GetEvent(ctx, eventID)
}
