package api

import (
	"event-schedule/internal/lib/logger/sl"
	"event-schedule/internal/model"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type GetVacantRoomsResponse struct {
	Response *Response
	Rooms    []*model.Suite
}

func (i *Implementation) GetVacantRooms(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.api.GetVacantRooms"
		var dates string
		var err error

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if dates = chi.URLParam(r, "interval"); dates == "" {
			log.Error("invalid request", sl.Err(ErrNoInterval))
			render.Render(w, r, ErrInvalidRequest(ErrNoInterval))
			return
		}

		interval := strings.Split(dates, ",")
		startDate, err := time.Parse("2006-01-02T15:04:05Z07:00", interval[0])
		endDate, err := time.Parse("2006-01-02T15:04:05Z07:00", interval[1])
		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, ErrInvalidRequest(ErrInvalidDateFormat))
			return
		}

		rooms, err := i.Service.GetVacantRooms(r.Context(), startDate, endDate) //TODO:GetVacantRooms
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, ErrInternalError(err))
			return
		}

		log.Info("vacant room acquired", slog.Any("rooms", rooms))
		render.Status(r, http.StatusCreated)
		render.Render(w, r, GetVacantRoomsAPI(rooms))
	}
}

func GetVacantRoomsAPI(rooms []*model.Suite) *GetVacantRoomsResponse {
	return &GetVacantRoomsResponse{
		Response: OK(),
		Rooms:    rooms,
	}
}

func (rd *GetVacantRoomsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
