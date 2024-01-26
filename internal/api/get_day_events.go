package api

import (
	"errors"
	"event-schedule/internal/model"
	"fmt"
	"log/slog"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type GetDayEventsResponse struct {
	Response *Response `json:"response"`
	//TODO: implement convert Event structs
	Events []*model.Event `json:"events"`
}

func (i *Implementation) GetDayEvents(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.api.GetDayEvents"
		var userID string

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		if userID = chi.URLParam(r, "userID"); userID == "" {
			render.Render(w, r, ErrInvalidRequest(errors.New("received no userID")))
			return
		}

		id, err := i.Service.GetDayEvents(r.Context(), userID)
		if err != nil {
			render.Render(w, r, ErrInternalError(err))
		}
		render.Status(r, http.StatusCreated)
		render.Render(w, r, GetDayEventsResponseAPI(fmt.Sprintf("received GetDayEvent from %s", id)))
	}

}

func GetDayEventsResponseAPI(eventIDs string) *GetDayEventsResponse {
	resp := &GetDayEventsResponse{
		Response: OK(),
		Events:   []*model.Event{{Uuid: eventIDs}, {Uuid: eventIDs + "1"}},
	}

	return resp
}

func (rd *GetDayEventsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
