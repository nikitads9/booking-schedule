package handlers

import (
	"errors"
	"event-schedule/internal/app/api"
	"event-schedule/internal/app/convert"
	"event-schedule/internal/logger/sl"
	"log/slog"

	"net/http"

	validator "github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// SignUp godoc
//
//	@Summary		Sign up
//	@Description	Creates user with given tg id, nickname, name and password hashed by bcrypto. Every parameter is required. Returns jwt token.
//	@ID				signUpUserJson
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param          user	body	api.User	true	"User"
//	@Success		200	{object}	api.AuthResponse
//	@Failure		400	{object}	api.AuthResponse
//	@Failure		404	{object}	api.AuthResponse
//	@Failure		503	{object}	api.AuthResponse
//	@Router			/user/sign-up [post]
func (i *Implementation) SignUp(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "user.api.handlers.SignUp"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		req := &api.User{}
		err := render.Bind(r, req)
		if err != nil {
			if errors.As(err, api.ValidateErr) {
				validateErr := err.(validator.ValidationErrors)
				log.Error("some of the required values were not received or were null", sl.Err(validateErr))
				render.Render(w, r, api.ErrValidationError(validateErr))
				return
			}
			log.Error("failed to decode request body", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}
		//TODO: remove
		log.Debug("request body decoded", slog.Any("req", req))
		user, err := convert.ToUserInfo(req)

		if err != nil {
			log.Error("invalid request", sl.Err(err))
			render.Render(w, r, api.ErrInvalidRequest(err))
			return
		}

		token, err := i.User.SignUp(ctx, user)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			render.Render(w, r, api.ErrInternalError(err))
			return
		}

		log.Debug("user created", slog.Any("token:", token))

		render.Status(r, http.StatusCreated)
		render.Render(w, r, api.AuthResponseAPI(token))
	}

}
