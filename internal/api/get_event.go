package api

import (
	"event-schedule/internal/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type GetEventResponse struct {
	Response *Response    `json:"response"`
	Event    *model.Event `json:"event"`
}

func (i *Implementation) GetEvent(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: проверить
		// Assume if we've reach this far, we can access the event
		// context because this handler is a child of the EventCtx
		// middleware. The worst case, the recoverer middleware will save us.
		event := r.Context().Value("event").(*model.Event)

		if err := render.Render(w, r, GetEventResponseAPI(event)); err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}
	}
}

func GetEventResponseAPI(event *model.Event) *GetEventResponse {
	resp := &GetEventResponse{
		Response: OK(),
		Event:    event,
	}

	return resp
}

func (rd *GetEventResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
