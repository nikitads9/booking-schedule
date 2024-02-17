package event

import (
	"context"
	"errors"
	"event-schedule/internal/client/db"
	"event-schedule/internal/model"
	t "event-schedule/internal/repository/table"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/go-chi/chi/middleware"
)

/*
	WITH booked_intervals AS (
	   SELECT
	       end_date AS booked_start,
	       LEAD(start_date) OVER (PARTITION BY suite_id ORDER BY start_date) AS booked_end
	   FROM events
	   WHERE suite_id = 1

),
free_intervals AS (

		SELECT
		    CASE
		        WHEN booked_start IS NULL THEN 0
		        ELSE booked_start
		    END AS start,
		    CASE
		        WHEN booked_end IS NULL THEN 30
		        ELSE booked_end
		    END AS end
		FROM booked_intervals
		UNION ALL
	    SELECT 0 AS start, 30 AS end
	    WHERE NOT EXISTS (SELECT 1 FROM booked_intervals)

)
SELECT * FROM free_intervals
WHERE start < 30;

по итогу вот такой запрос
"WITH booked_intervals AS ( SELECT end_date AS booked_start, LEAD(start_date) OVER (PARTITION BY suite_id ORDER BY start_date) AS booked_end FROM events WHERE suite_id = $1 ) SELECT * FROM (SELECT (CASE WHEN booked_start is NULL THEN 02-07-2024 ELSE booked_start END) AS start, (CASE WHEN booked_end IS NULL THEN 03-08-2024 ELSE booked_end END) AS end FROM booked_intervals UNION ALL SELECT 02-07-2024 AS start , 03-08-2024 AS end  WHERE NOT EXISTS (SELECT 1 FROM booked_intervals)) AS free_intervals WHERE start < $1"
*/
func (r *repository) GetVacantDates(ctx context.Context, suiteID int64) ([]*model.Interval, error) {
	const op = "events.repository.GetVacantDates"

	now := time.Now().Format("01-02-2006")
	month := time.Now().Add(720 * time.Hour).Format("01-02-2006")

	r.log = r.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	bookedQuery, bookedArgs, err := sq.Select(t.EndDate+" AS booked_start", "LEAD("+t.StartDate+") OVER (PARTITION BY "+t.SuiteID+" ORDER BY "+t.StartDate+") AS booked_end").
		From(t.EventTable).
		Where(sq.Eq{t.SuiteID: suiteID}).
		Prefix("WITH booked_intervals AS (").
		Suffix(")").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		r.log.Error("failed to build a booked_intervals subquery", err)
		return nil, ErrQueryBuild
	}

	noEventsQuery, noEventsArgs, err := sq.Select(now+" AS start ", month+" AS end ").
		Where("NOT EXISTS (SELECT 1 FROM booked_intervals)").
		ToSql()
	if err != nil {
		r.log.Error("failed to build check on no bookings subquery", err)
		return nil, ErrQueryBuild
	}

	start := sq.Case().
		When("booked_start is NULL", now).
		Else("booked_start")
	end := sq.Case().
		When("booked_end IS NULL", month).
		Else("booked_end")
	free := sq.Select().
		Column(sq.Alias(start, "start")).
		Column(sq.Alias(end, "end")).
		From("booked_intervals").
		Suffix("UNION ALL "+noEventsQuery, noEventsArgs...)

	builder := sq.Select("*").
		FromSelect(free, "free_intervals").
		Where(sq.Lt{"start": month}).
		Prefix(bookedQuery, bookedArgs...).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		r.log.Error("failed to build a query", err)
		return nil, ErrQueryBuild
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var res []*model.Interval
	err = r.client.DB().SelectContext(ctx, res, q, args...)
	if err != nil {
		if errors.As(err, pgNoConnection) {
			r.log.Error("no connection to database host", err)
			return nil, ErrNoConnection
		}
		if pgxscan.NotFound(err) {
			r.log.Error("no vacant dates within month for this room", err)
			return nil, ErrNotFound
		}
		r.log.Error("query execution error", err)
		return nil, ErrQuery
	}

	return res, nil
}
