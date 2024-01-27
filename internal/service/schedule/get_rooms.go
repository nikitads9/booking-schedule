package schedule

import (
	"context"
	"event-schedule/internal/model"

	gofakeit "github.com/brianvoe/gofakeit/v6"

	"time"
)

func (s *Service) GetVacantRooms(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.Suite, error) {
	return []*model.Suite{
		{
			SuiteID:  gofakeit.Int64(),
			Capacity: gofakeit.Int8(),
			Name:     gofakeit.City(),
		},
	}, nil
}
