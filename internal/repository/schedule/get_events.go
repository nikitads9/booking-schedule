package schedule

import (
	"context"
	"event-schedule/internal/model"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofrs/uuid"
)

func (r *repository) GetEvents(ctx context.Context, mod *model.GetEventsInfo) ([]*model.EventInfo, error) {

	//builder := sq.Select("").PlaceholderFormat(sq.Dollar)
	return []*model.EventInfo{
		{EventID: uuid.Nil,
			SuiteID:   gofakeit.Int64(),
			StartDate: gofakeit.FutureDate(),
			EndDate:   gofakeit.FutureDate(),
			CreatedAt: gofakeit.PastDate()},
		{EventID: uuid.Nil,
			SuiteID:   gofakeit.Int64(),
			StartDate: gofakeit.FutureDate(),
			EndDate:   gofakeit.FutureDate(),
		},
	}, nil
}
