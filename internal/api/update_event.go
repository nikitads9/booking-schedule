package api

import (
	"event-schedule/internal/lib/logger/sl"
	"event-schedule/internal/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type UpdateEventRequest struct {
	EventInfo *model.EventInfo `json:"eventInfo"`
}

type UpdateEventResponse struct {
	Response *Response `json:"response"`
}

// UpdateEvent updates an existing Event in our persistent store.
func (i *Implementation) UpdateEvent(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.events.api.UpdateEvent"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		event := r.Context().Value("event").(*model.Event)

		//TODO: getter method getEventInfo
		data := &UpdateEventRequest{EventInfo: event.EventInfo}
		if err := render.Bind(r, data); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		log.Info("request body decoded", slog.Any("req", data))

		//dbUpdateEvent(event.Uuid, data.EventInfo)
		//обработка ошибки

		log.Info("event updated", slog.Any("id", event.UUID))

		render.Render(w, r, UpdateEventResponseAPI())
	}
}

func (e *UpdateEventRequest) Bind(r *http.Request) error {
	return nil
}

func (rd *UpdateEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func UpdateEventResponseAPI() *UpdateEventResponse {
	return &UpdateEventResponse{
		Response: OK(),
	}
}
