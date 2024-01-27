package schedule

import (
	"context"
	"event-schedule/internal/model"

	"github.com/brianvoe/gofakeit/v6"
)

func (r *repository) GetEvents(ctx context.Context, userID string) ([]*model.Event, error) {

	return []*model.Event{
		{UUID: gofakeit.UUID(),
			EventInfo: &model.EventInfo{
				SuiteID:   gofakeit.Int64(),
				SuiteName: gofakeit.City(),
				StartDate: gofakeit.FutureDate(),
				EndDate:   gofakeit.FutureDate(),
				OwnerID:   gofakeit.Int64(),
			},
			CreatedAt: gofakeit.PastDate()},
		{UUID: gofakeit.UUID(),
			EventInfo: &model.EventInfo{
				SuiteID:   gofakeit.Int64(),
				SuiteName: gofakeit.City(),
				StartDate: gofakeit.FutureDate(),
				EndDate:   gofakeit.FutureDate(),
				OwnerID:   gofakeit.Int64(),
			},
		}}, nil
}
