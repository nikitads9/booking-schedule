package schedule

import (
	"context"
	"event-schedule/internal/model"

	"time"
)

func (s *Service) GetVacantDates(ctx context.Context, suiteID int64) ([]*model.Interval, error) {
	return []*model.Interval{
		{
			StartDate: time.Now(),
			EndDate:   time.Now().Add(time.Hour),
		},
	}, nil
}
