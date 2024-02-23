package event

import (
	"context"

	"github.com/gofrs/uuid"
)

func (s *Service) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	return s.eventRepository.DeleteEvent(ctx, eventID)
}
