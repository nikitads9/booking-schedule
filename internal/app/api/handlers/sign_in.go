package handlers

import (
	"errors"
	"event-schedule/internal/app/api"
	"event-schedule/internal/app/service/user"
	"event-schedule/internal/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// SignIn godoc
//
//	@Summary		Sign in
//	@Description	Get auth token to access user restricted api methods. Requires nickname and password passed via basic auth.
//	@ID				getOauthToken
//	@Tags			users
//	@Produce		json
//	@Success		200	{object}	api.AuthResponse
//	@Failure		400	{object}	api.AuthResponse
//	@Failure		401	{object}	api.AuthResponse
//	@Failure		503	{object}	api.AuthResponse
//	@Router			/user/sign-in [get]
//
// @Security 		BasicAuth
func (i *Implementation) SignIn(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "events.api.auth.SignIn"

		ctx := r.Context()
		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		nickname, pass, ok := r.BasicAuth()
		if !ok {
			log.Error("bad request")
			render.Render(w, r, api.ErrInvalidRequest(api.ErrNoAuth))
			return
		}
		token, err := i.User.SignIn(ctx, nickname, pass)
		if err != nil {
			if errors.Is(err, user.ErrBadLogin) {
				log.Error("incorrect login", sl.Err(err))
				render.Render(w, r, api.ErrUnauthorized(err))
				return
			}
			if errors.Is(err, user.ErrBadPasswd) {
				log.Error("incorrect passwd", sl.Err(err))
				render.Render(w, r, api.ErrUnauthorized(err))
				return
			}
			log.Error("failed to login user: ", err)
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		render.Render(w, r, api.AuthResponseAPI(token))
	}

}
