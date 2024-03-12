package booking

import (
	"booking-schedule/internal/app/model"
	"context"

	"github.com/gofrs/uuid"
)

func (s *Service) GetBooking(ctx context.Context, bookingID uuid.UUID, userID int64) (*model.BookingInfo, error) {
	return s.bookingRepository.GetBooking(ctx, bookingID, userID)
}
