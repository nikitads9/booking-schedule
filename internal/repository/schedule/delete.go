package schedule

import (
	"context"
	"event-schedule/internal/client/db"
	"event-schedule/internal/repository/table"

	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
)

func (r *repository) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	builder := sq.Delete(table.EventTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "event_repository.DeleteEvent",
		QueryRaw: query,
	}

	result, err := r.client.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
