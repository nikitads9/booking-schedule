package event

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"event-schedule/internal/client/db"
	"event-schedule/internal/model"
	t "event-schedule/internal/repository/table"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-chi/chi/middleware"
)

//TODO: проверить а что будет еси комнаты забронированы гетерогенно: разными клиентами и приходит запрос на обновление одним из них, накладывающийся на второго (умозрительно вроже норм)
/*
SELECT NOT EXISTS ( SELECT 1 FROM events WHERE ((suite_id = 1 AND ((start_date > 3 AND start_date < 6) OR (end_date > 3 AND end_date < 6)))) ) as availible, (SELECT EXISTS ( SELECT 1 FROM events WHERE ((suite_id = 1 AND owner_id = 2) AND ((start_date > 3 AND start_date < 6) OR (end_date > 3 AND end_date < 6)))) ) as occupied_by_owner
*/
func (r *repository) CheckAvailibility(ctx context.Context, suiteID int64, start time.Time, end time.Time, userID int64) (*model.Availibility, error) {
	const op = "events.repository.CheckAvailibility"

	r.log = r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	subQuery := sq.Select("1").From(t.EventTable).Where(sq.And{
		sq.And{
			sq.And{
				sq.Eq{t.SuiteID: suiteID},
				sq.Eq{t.OwnerID: userID},
			},
			sq.Or{
				sq.And{
					sq.Gt{t.StartDate: start},
					sq.Lt{t.StartDate: end},
				},
				sq.And{
					sq.Gt{t.EndDate: start},
					sq.Lt{t.EndDate: end},
				},
			},
		},
	}).
		Prefix("(SELECT EXISTS (").
		Suffix(")) as occupied_by_owner").
		PlaceholderFormat(sq.Dollar)

	query, args, err := sq.Select("1").From(t.EventTable).Where(sq.And{
		sq.And{
			sq.Eq{t.SuiteID: suiteID},
			sq.Or{
				sq.And{
					sq.Gt{t.StartDate: start},
					sq.Lt{t.StartDate: end},
				},
				sq.And{
					sq.Gt{t.EndDate: start},
					sq.Lt{t.EndDate: end},
				},
			},
		},
	}).
		Prefix("SELECT NOT EXISTS (").
		Suffix(") as availible,").
		SuffixExpr(subQuery).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		r.log.Error("failed to build query", err)
		return nil, ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}
	//TODO: ensure that model is parsed correctly
	var res = new(model.Availibility)
	err = r.client.DB().SelectContext(ctx, &res, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			r.log.Error("no connection to database host", err)
			return nil, ErrNoConnection
		}
		r.log.Error("query execution error", err)
		return nil, ErrQuery
	}

	return res, nil
}
