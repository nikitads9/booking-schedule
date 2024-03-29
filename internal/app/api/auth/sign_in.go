package auth

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
//	@Failure		400	{object}	api.errResponse
//	@Failure		401	{object}	api.errResponse
//	@Failure		503	{object}	api.errResponse
//	@Router			/sign-in [get]
//
//	@Security 		BasicAuth
func (i *Implementation) SignIn(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.auth.SignIn"

		ctx := r.Context()
		requestID := middleware.GetReqID(ctx)

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", requestID),
		)
		ctx, span := i.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
		defer span.End()

		nickname, pass, ok := r.BasicAuth()
		if !ok {
			span.RecordError(api.ErrBadRequest)
			span.SetStatus(codes.Error, api.ErrBadRequest.Error())
			log.Error("bad request")
			api.WriteWithError(w, http.StatusBadRequest, api.ErrNoAuth.Error())
			return
		}

		span.AddEvent("acquired login and password")

		token, err := i.user.SignIn(ctx, nickname, pass)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to sign in user", sl.Err(err))
			api.WriteWithError(w, GetErrorCode(err), err.Error())
			return
		}

		span.AddEvent("signed token acquired")
		log.Info("user signed in", slog.Any("login", nickname))

		api.WriteWithStatus(w, http.StatusOK, api.AuthResponse{
			Token: token,
		})
	}

}
