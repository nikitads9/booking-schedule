package api

import (
	"event-schedule/internal/lib/logger/sl"
	"event-schedule/internal/model"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type GetVacantDatesResponse struct {
	Response  *Response
	Intervals []*model.Interval
}

func (i *Implementation) GetVacantDates(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.api.GetVacantDates"
		var suiteID string
		var id int64
		var err error

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if suiteID = chi.URLParam(r, "suiteID"); suiteID == "" {
			log.Error("invalid request", sl.Err(ErrNoSuiteID))
			render.Render(w, r, ErrInvalidRequest(ErrNoSuiteID))
			return
		}

		if id, err = strconv.ParseInt(suiteID, 10, 64); err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		intervals, err := i.Service.GetVacantDates(r.Context(), id) //TODO:GetVacantDates
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, ErrInternalError(err))
			return
		}

		log.Info("vacant dates acquired", slog.Any("intervals", intervals))

		render.Status(r, http.StatusCreated)
		render.Render(w, r, GetVacantDatesAPI(intervals))
	}
}

func GetVacantDatesAPI(intervals []*model.Interval) *GetVacantDatesResponse {
	return &GetVacantDatesResponse{
		Response:  OK(),
		Intervals: intervals,
	}
}

func (rd *GetVacantDatesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
