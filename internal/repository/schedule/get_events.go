package schedule

import (
	"context"
	"event-schedule/internal/model"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofrs/uuid"
)

func (r *repository) GetEvents(ctx context.Context, userID int64) ([]*model.EventInfo, error) {
	uid := gofakeit.UUID()
	res, _ := uuid.FromString(uid)
	return []*model.EventInfo{
		{EventID: res,
			SuiteID:   gofakeit.Int64(),
			StartDate: gofakeit.FutureDate(),
			EndDate:   gofakeit.FutureDate(),
			CreatedAt: gofakeit.PastDate()},
		{EventID: res,
			SuiteID:   gofakeit.Int64(),
			StartDate: gofakeit.FutureDate(),
			EndDate:   gofakeit.FutureDate(),
		},
	}, nil
}
