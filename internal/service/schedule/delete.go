package schedule

import (
	"context"

	"github.com/gofrs/uuid"
)

func (s *Service) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	return s.scheduleRepository.DeleteEvent(ctx, eventID)
}
