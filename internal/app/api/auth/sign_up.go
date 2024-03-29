package auth

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
	"errors"
	"log/slog"

	"net/http"

	validator "github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

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
//	@Failure		400	{object}	api.errResponse
//	@Failure		503	{object}	api.errResponse
//	@Router			/sign-up [post]
func (i *Implementation) SignUp(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.auth.SignUp"

		ctx := r.Context()
		requestID := middleware.GetReqID(ctx)

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", requestID),
		)
		ctx, span := i.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
		defer span.End()

		req := &api.SignUpRequest{}
		err := render.Bind(r, req)
		if err != nil {
			if errors.As(err, api.ValidateErr) {
				validateErr := err.(validator.ValidationErrors)
				span.RecordError(validateErr)
				span.SetStatus(codes.Error, validateErr.Error())
				log.Error("some of the required values were not received or were null", sl.Err(validateErr))
				api.WriteValidationError(w, validateErr)
				return
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to decode request body", sl.Err(err))
			api.WriteWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		span.AddEvent("request body decoded")

		user, err := convert.ToUserInfo(req)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("invalid request", sl.Err(err))
			api.WriteWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		span.AddEvent("request model converted")

		token, err := i.user.SignUp(ctx, user)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("sign up failed", sl.Err(err))
			api.WriteWithError(w, GetErrorCode(err), err.Error())
			return
		}

		span.AddEvent("user created")
		log.Info("user created", slog.Any("login: ", req.Nickname))

		render.Status(r, http.StatusCreated)
		api.WriteWithStatus(w, http.StatusOK, api.AuthResponse{
			Token: token,
		})
	}

}
