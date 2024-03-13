package auth

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
	"errors"
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
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param          user	body	api.SignUpRequest	true	"User"
//	@Success		200	{object}	api.AuthResponse
//	@Failure		400	{object}	api.AuthResponse
//	@Failure		404	{object}	api.AuthResponse
//	@Failure		503	{object}	api.AuthResponse
//	@Router			/auth/sign-up [post]
func (i *Implementation) SignUp(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "user.SignUp"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		_, span := i.Tracer.Start(ctx, op)
		defer span.End()

		req := &api.SignUpRequest{}
		err := render.Bind(r, req)
		if err != nil {
			if errors.As(err, api.ValidateErr) {
				validateErr := err.(validator.ValidationErrors)
				log.Error("some of the required values were not received or were null", sl.Err(validateErr))
				err = render.Render(w, r, api.ErrValidationError(validateErr))
				if err != nil {
					log.Error("failed to render response", sl.Err(err))
					return
				}
				return
			}
			log.Error("failed to decode request body", sl.Err(err))
			err = render.Render(w, r, api.ErrInvalidRequest(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		user, err := convert.ToUserInfo(req)

		if err != nil {
			log.Error("invalid request", sl.Err(err))
			err = render.Render(w, r, api.ErrInvalidRequest(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		token, err := i.User.SignUp(ctx, user)
		if err != nil {
			log.Error("internal error", sl.Err(err))
			err = render.Render(w, r, api.ErrInternalError(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		log.Debug("user created", slog.Any("token: ", token))
		log.Info("user created", slog.Any("login: ", req.Nickname))

		render.Status(r, http.StatusCreated)
		err = render.Render(w, r, api.AuthResponseAPI(token))
		if err != nil {
			log.Error("failed to render response", sl.Err(err))
			return
		}
	}

}
