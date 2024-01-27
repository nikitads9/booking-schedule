package schedule

import (
	"context"
	"event-schedule/internal/model"

	gofakeit "github.com/brianvoe/gofakeit/v6"
)

func (r *repository) GetEvent(ctx context.Context, eventID string) (*model.Event, error) {

	return &model.Event{
		UUID: eventID,
		EventInfo: &model.EventInfo{
			SuiteID:   gofakeit.Int64(),
			StartDate: gofakeit.FutureDate(),
			EndDate:   gofakeit.FutureDate(),
			OwnerID:   gofakeit.Int64(),
		},
		CreatedAt: gofakeit.PastDate(),
	}, nil
}
