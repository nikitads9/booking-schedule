package user

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/middleware/auth"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/codes"
)

// GetMyProfile godoc
//
//	@Summary		Get info for current user
//	@Description	Responds with account info for signed in user.
//	@ID				getMyUserAuth
//	@Tags			users
//	@Produce		json
//
//	@Success		200	{object}	api.GetMyProfileResponse
//	@Failure		400	{object}	api.errResponse
//	@Failure		401	{object}	api.errResponse
//	@Failure		404	{object}	api.errResponse
//	@Failure		503	{object}	api.errResponse
//	@Router			/user/me [get]
//
// @Security Bearer
func (i *Implementation) GetMyProfile(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.user.GetMyProfile"

		ctx := r.Context()

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(ctx)),
		)
		ctx, span := i.tracer.Start(ctx, op)
		defer span.End()

		userID := auth.UserIDFromContext(ctx)
		if userID == 0 {
			span.RecordError(api.ErrNoUserID)
			span.SetStatus(codes.Error, api.ErrNoUserID.Error())
			log.Error("no user id in context", sl.Err(api.ErrNoUserID))
			api.WriteWithError(w, http.StatusUnauthorized, api.ErrNoAuth.Error())
			return
		}

		user, err := i.user.GetUser(ctx, userID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("internal error", sl.Err(err))
			api.WriteWithError(w, GetErrorCode(err), err.Error())
			return
		}

		span.AddEvent("user acquired")
		log.Info("user acquired", slog.Any("user: ", user))

		api.WriteWithStatus(w, http.StatusOK, api.GetMyProfileResponse{
			Profile: convert.ToApiUserInfo(user),
		})
	}
}
