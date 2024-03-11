package handlers

import (
	"event-schedule/internal/app/api"
	"event-schedule/internal/app/convert"
	"event-schedule/internal/logger/sl"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// GetVacantRooms godoc
//
//	@Summary		Get list of vacant rooms
//	@Description	Receives two dates as query parameters. start is to be before end and both should not be expired. Responds with list of vacant rooms and their parameters for given interval.
//	@ID				getRoomsByDates
//	@Tags			bookings
//	@Produce		json
//	@Param			start	query	string	true	"start"	Format(time.Time) default(2024-03-28T17:43:00Z)
//	@Param			end	query	string	true	"end"	Format(time.Time) default(2024-03-29T17:43:00Z)
//	@Success		200	{object}	api.GetVacantRoomsResponse
//	@Failure		400	{object}	api.GetVacantRoomsResponse
//	@Failure		404	{object}	api.GetVacantRoomsResponse
//	@Failure		422	{object}	api.GetVacantRoomsResponse
//	@Failure		503	{object}	api.GetVacantRoomsResponse
//	@Router			/get-vacant-rooms [get]
func (i *Implementation) GetVacantRooms(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "events.api.handlers.GetVacantRooms"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		start := r.URL.Query().Get("start")
		if start == "" {
			log.Error("invalid request", sl.Err(api.ErrNoInterval))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoInterval))
			return
		}

		end := r.URL.Query().Get("end")
		if end == "" {
			log.Error("invalid request", sl.Err(api.ErrNoInterval))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoInterval))
			return
		}

		startDate, err := time.Parse("2006-01-02T15:04:05Z", start)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrParse))
			return
		}
		endDate, err := time.Parse("2006-01-02T15:04:05Z", end)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(api.ErrParse))
			return
		}

		err = api.CheckDates(startDate, endDate)
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
		}

		rooms, err := i.Booking.GetVacantRooms(ctx, startDate, endDate)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Info("vacant rooms acquired", slog.Any("quantity:", len(rooms)))
		render.Status(r, http.StatusCreated)
		render.Render(w, r, api.GetVacantRoomsAPI(convert.ToApiSuites(rooms)))
	}
}
