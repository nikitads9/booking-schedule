package booking

import (
	"booking-schedule/internal/app/model"
	"context"
	"time"
)

func (s *Service) GetVacantRooms(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.Suite, error) {
	return s.bookingRepository.GetVacantRooms(ctx, startDate, endDate)
}
