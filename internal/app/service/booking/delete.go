package booking

import (
	"context"

	"github.com/gofrs/uuid"
)

func (s *Service) DeleteBooking(ctx context.Context, bookingID uuid.UUID, userID int64) error {
	return s.bookingRepository.DeleteBooking(ctx, bookingID, userID)
}
