package event

import (
	"context"
	"errors"
	"log/slog"

	"event-schedule/internal/app/model"
	t "event-schedule/internal/app/repository/table"
	"event-schedule/internal/pkg/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-chi/chi/middleware"
)

//TODO: проверить а что будет еси комнаты забронированы гетерогенно: разными клиентами и приходит запрос на обновление одним из них, накладывающийся на второго (умозрительно вроже норм)
/*
SELECT NOT EXISTS ( SELECT 1 FROM events WHERE ((suite_id = 1 AND ((start_date > 3 AND start_date < 6) OR (end_date > 3 AND end_date < 6)))) ) as availible, (SELECT EXISTS ( SELECT 1 FROM events WHERE ((suite_id = 1 AND owner_id = 2) AND ((start_date > 3 AND start_date < 6) OR (end_date > 3 AND end_date < 6)))) ) as occupied_by_owner
*/
//TODO: проверять занято ли автором именно та бронь, которую он хочет отредактировать?
func (r *repository) CheckAvailibility(ctx context.Context, mod *model.Event) (*model.Availibility, error) {
	const op = "events.repository.CheckAvailibility"

	r.log = r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	subQuery := sq.Select("1").From(t.EventTable).Where(sq.And{
		sq.And{
			sq.And{
				sq.Eq{t.SuiteID: mod.SuiteID},
				sq.And{sq.Eq{t.OwnerID: mod.UserID},
					sq.Eq{t.ID: mod.GetEventID()},
				},
			},
			sq.Or{
				sq.And{
					sq.Gt{t.StartDate: mod.StartDate},
					sq.Lt{t.StartDate: mod.EndDate},
				},
				sq.And{
					sq.Gt{t.EndDate: mod.StartDate},
					sq.Lt{t.EndDate: mod.EndDate},
				},
			},
		},
	}).
		Prefix("(SELECT EXISTS (").
		Suffix(")) as occupied_by_client").
		PlaceholderFormat(sq.Dollar)

	query, args, err := sq.Select("1").From(t.EventTable).Where(sq.And{
		sq.And{
			sq.Eq{t.SuiteID: mod.SuiteID},
			sq.Or{
				sq.And{
					sq.Gt{t.StartDate: mod.StartDate},
					sq.Lt{t.StartDate: mod.EndDate},
				},
				sq.And{
					sq.Gt{t.EndDate: mod.StartDate},
					sq.Lt{t.EndDate: mod.EndDate},
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
	err = r.client.DB().GetContext(ctx, res, q, args...)
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
