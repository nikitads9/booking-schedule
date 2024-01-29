package schedule

import (
	"context"
	"event-schedule/internal/client/db"
	"event-schedule/internal/model"
	"event-schedule/internal/repository/table"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (r *repository) UpdateEvent(ctx context.Context, mod *model.UpdateEventInfo) error {
	//TODO: check this at service level
	if !mod.SuiteID.Valid && !mod.StartDate.Valid && !mod.EndDate.Valid && !mod.NotificationPeriod.Valid {
		return ErrFailed
	}

	builder := sq.Update(table.EventTable).
		Set("updated_at", time.Now().UTC()).
		Where(sq.Eq{"id": mod.EventID}).
		PlaceholderFormat(sq.Dollar)

	if mod.SuiteID.Valid {
		builder.Set("suite_id", mod.SuiteID.Int64)
	}

	if mod.StartDate.Valid {
		builder.Set("start_date", mod.StartDate.Time)
	}

	if mod.EndDate.Valid {
		builder.Set("end_date", mod.EndDate.Time)
	}

	if mod.NotificationPeriod.Valid {
		notificationPeriod, err := time.ParseDuration(mod.NotificationPeriod.String)
		if err != nil {
			return ErrFailed
		}
		builder.Set("notification_period", notificationPeriod)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "event_repository.UpdateEvent",
		QueryRaw: query,
	}

	_, err = r.client.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}
