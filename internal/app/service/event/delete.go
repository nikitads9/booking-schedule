package event

import (
	"context"

	"github.com/gofrs/uuid"
)

func (s *Service) DeleteEvent(ctx context.Context, eventID uuid.UUID, userID int64) error {
	return s.eventRepository.DeleteEvent(ctx, eventID, userID)
}
