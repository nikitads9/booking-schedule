package auth

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
	"errors"
	"log/slog"

	"net/http"

	validator "github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/codes"

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
//	@Router			/sign-up [post]
func (i *Implementation) SignUp(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.auth.SignUp"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)

		ctx, span := i.tracer.Start(ctx, op)
		defer span.End()

		req := &api.SignUpRequest{}
		err := render.Bind(r, req)
		if err != nil {
			if errors.As(err, api.ValidateErr) {
				validateErr := err.(validator.ValidationErrors)
				span.RecordError(validateErr)
				span.SetStatus(codes.Error, validateErr.Error())
				log.Error("some of the required values were not received or were null", sl.Err(validateErr))
				err = render.Render(w, r, api.ErrValidationError(validateErr))
				if err != nil {
					log.Error("failed to render response", sl.Err(err))
					return
				}
				return
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to decode request body", sl.Err(err))
			err = render.Render(w, r, api.ErrInvalidRequest(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("request body decoded")

		user, err := convert.ToUserInfo(req)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("invalid request", sl.Err(err))
			err = render.Render(w, r, api.ErrInvalidRequest(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("request model converted")

		token, err := i.user.SignUp(ctx, user)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("internal error", sl.Err(err))
			err = render.Render(w, r, api.ErrInternalError(err))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("user created")
		log.Info("user created", slog.Any("login: ", req.Nickname))

		render.Status(r, http.StatusCreated)
		err = render.Render(w, r, api.AuthResponseAPI(token))
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to render response", sl.Err(err))
			return
		}
	}

}
