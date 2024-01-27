package api

import (
	"event-schedule/internal/lib/logger/sl"
	"event-schedule/internal/model"
	"log/slog"
	"strconv"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type GetEventsResponse struct {
	Response *Response `json:"response"`
	//TODO: implement convert Event structs
	Events []*model.Event `json:"events"`
}

func (i *Implementation) GetEvents(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.api.GetEvents"
		var userID, period string
		var id int64
		var err error

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if userID = chi.URLParam(r, "userID"); userID == "" {
			log.Error("invalid request", sl.Err(ErrNoUserID))
			render.Render(w, r, ErrInvalidRequest(ErrNoUserID))
			return
		}

		if id, err = strconv.ParseInt(userID, 10, 64); err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		if period = chi.URLParam(r, "period"); period == "" {
			log.Error("invalid request", sl.Err(ErrNoInterval))
			render.Render(w, r, ErrInvalidRequest(ErrNoInterval))
			return
		}

		events, err := i.Service.GetEvents(r.Context(), id, period) //TODO: filter
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, ErrInternalError(err))
			return
		}

		log.Info("events acquired", slog.Int("quantity", len(events)))

		render.Status(r, http.StatusCreated)
		render.Render(w, r, GetEventsResponseAPI(events))
	}

}

func GetEventsResponseAPI(events []*model.Event) *GetEventsResponse {
	resp := &GetEventsResponse{
		Response: OK(),
		Events:   events,
	}

	return resp
}

func (rd *GetEventsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
