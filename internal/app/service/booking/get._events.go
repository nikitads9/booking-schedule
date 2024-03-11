package booking

import (
	"booking-schedule/internal/app/model"
	"context"
	"time"
)

func (s *Service) GetBookings(ctx context.Context, startDate time.Time, endDate time.Time, id int64) ([]*model.BookingInfo, error) {
	return s.bookingRepository.GetBookings(ctx, startDate, endDate, id)
}
