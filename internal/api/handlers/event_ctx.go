package handlers

// TODO: remove EventCTX
// EventCtx middleware is used to load an Event object from
// the URL parameters passed through as the request. In case
// the Event could not be found, we stop here and return a 404.
/* func (i *Implementation) EventCtx(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "handlers.events.api.EventCtx"

			ctx := r.Context()

			log = log.With(
				slog.String("op", op),
				slog.String("request_id", middleware.GetReqID(ctx)),
			)

			eventID := chi.URLParam(r, "event_id")
			if eventID == "" {
				log.Error("invalid request", sl.Err(api.ErrNoEventID))
				render.Render(w, r, api.ErrInvalidRequest(api.ErrNoEventID))
				return
			}

			eventUUID, err := uuid.FromString(eventID)
			if err != nil {
				log.Error("invalid request", sl.Err(err))
				render.Render(w, r, api.ErrInvalidRequest(api.ErrParse))
				return
			}

			if eventUUID == uuid.Nil {
				log.Error("invalid request", sl.Err(api.ErrNoEventID))
				render.Render(w, r, api.ErrInvalidRequest(api.ErrNoEventID))
				return
			}

			log.Info("decoded URL param", slog.Any("eventID", eventUUID))

			 			event, err = i.Service.GetEvent(ctx, eventUUID)
			   			if err != nil {
			   				log.Error("internal error", sl.Err(err))
			   				render.Render(w, r, api.ErrInternalError(err))
			   				return
			   			}

			   			log.Info("event acquired", slog.Any("event", event))

			ctx = context.WithValue(ctx, "event", eventUUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
} */
