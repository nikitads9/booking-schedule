package user

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/middleware/auth"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// DeleteMyProfile godoc
//
//	@Summary		Delete my profile
//	@Description	Deletes user and all bookings associated with him
//	@ID				deleteMyInfo
//	@Tags			users
//	@Produce		json
//
//	@Success		200	{object}	api.DeleteMyProfileResponse
//	@Failure		400	{object}	api.DeleteMyProfileResponse
//	@Failure		401	{object}	api.DeleteMyProfileResponse
//	@Failure		404	{object}	api.DeleteMyProfileResponse
//	@Failure		422	{object}	api.DeleteMyProfileResponse
//	@Failure		503	{object}	api.DeleteMyProfileResponse
//	@Router			/user/delete [delete]
//
// @Security Bearer
func (i *Implementation) DeleteMyProfile(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.user.DeleteMyProfile"

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
			err := render.Render(w, r, api.ErrUnauthorized(api.ErrNoAuth))
			if err != nil {
				log.Error("failed to render response", sl.Err(err))
				return
			}
			return
		}

		span.AddEvent("userID extracted from context", trace.WithAttributes(attribute.Int64("id", userID)))

		err := i.user.DeleteUser(ctx, userID)
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

		span.AddEvent("user deleted", trace.WithAttributes(attribute.Int64("id", userID)))
		log.Info("deleted booking", slog.Int64("id: ", userID))

		err = render.Render(w, r, api.DeleteMyProfileResponseAPI())
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to render response", sl.Err(err))
			return
		}
	}
}
