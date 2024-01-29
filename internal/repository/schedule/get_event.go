package schedule

import (
	"context"
	"event-schedule/internal/client/db"
	"event-schedule/internal/model"
	"event-schedule/internal/repository/table"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/gofrs/uuid"
)

func (r *repository) GetEvent(ctx context.Context, eventID uuid.UUID) (*model.EventInfo, error) {
	builder := sq.Select("id", "suite_id", "start_date", "end_date", "notification_period", "created_at", "updated_at").
		From(table.EventTable).
		Where(sq.Eq{"id": eventID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "event_repository.GetEvent",
		QueryRaw: query,
	}

	var res = new(model.EventInfo)
	err = r.client.DB().GetContext(ctx, res, q, args...)
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return res, nil
}
