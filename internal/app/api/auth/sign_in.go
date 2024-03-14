package auth

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/service/user"
	"booking-schedule/internal/logger/sl"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/codes"
)

// SignIn godoc
//
//	@Summary		Sign in
//	@Description	Get auth token to access user restricted api methods. Requires nickname and password passed via basic auth.
//	@ID				getOauthToken
//	@Tags			auth
//	@Produce		json
//
//	@Success		200	{object}	api.AuthResponse
//	@Failure		400	{object}	api.AuthResponse
//	@Failure		401	{object}	api.AuthResponse
//	@Failure		503	{object}	api.AuthResponse
//	@Router			/auth/sign-in [get]
//
//	@Security 		BasicAuth
func (i *Implementation) SignIn(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.auth.SignIn"

		ctx := r.Context()
		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		ctx, span := i.tracer.Start(ctx, op)
		defer span.End()

		nickname, pass, ok := r.BasicAuth()
		if !ok {
			span.RecordError(api.ErrBadRequest)
			span.SetStatus(codes.Error, api.ErrBadRequest.Error())
			log.Error("bad request")
			err := render.Render(w, r, api.ErrInvalidRequest(api.ErrNoAuth))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("acquired login and password")

		token, err := i.user.SignIn(ctx, nickname, pass)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			if errors.Is(err, user.ErrBadLogin) {
				log.Error("incorrect login", sl.Err(err))
				err = render.Render(w, r, api.ErrUnauthorized(err))
				if err != nil {
					log.Error("failed to render response", sl.Err(err))
					return
				}
				return
			}
			if errors.Is(err, user.ErrBadPasswd) {
				log.Error("incorrect passwd", sl.Err(err))
				err = render.Render(w, r, api.ErrUnauthorized(err))
				if err != nil {
					log.Error("failed to render response", sl.Err(err))
					return
				}
				return
			}
			log.Error("failed to login user: ", err)
			err = render.Render(w, r, api.ErrInternalError(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("token acquired")
		log.Info("user signed in", slog.Any("login: ", nickname))

		err = render.Render(w, r, api.AuthResponseAPI(token))
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to render response", sl.Err(err))
			return
		}
	}

}
