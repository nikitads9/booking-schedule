package schedule

import (
	"context"
	"event-schedule/internal/client/db"
	"event-schedule/internal/model"
	"event-schedule/internal/repository/table"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
)

func (r *repository) AddEvent(ctx context.Context, mod *model.Event) (uuid.UUID, error) {
	builder := sq.Insert(table.EventTable).
		Columns("user_id", "suite_id", "start_date", "end_date", "notification_period").
		Values(mod.UserID, mod.SuiteID, mod.StartDate, mod.EndDate, mod.NotificationPeriod).
		Suffix("returning id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	q := db.Query{
		Name:     "event_repository.AddEvent",
		QueryRaw: query,
	}

	row, err := r.client.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	row.Next()
	err = row.Scan(&id)
	if err != nil {
		return uuid.Nil, ErrFailed

	}

	return id, nil
}
