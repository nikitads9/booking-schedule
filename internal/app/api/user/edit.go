package user

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/convert"
	"booking-schedule/internal/logger/sl"
	"booking-schedule/internal/middleware/auth"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	validator "github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// EditMyProfile godoc
//
//	@Summary		Modify profile
//	@Description	Updates user's profile with provided values. If no values provided, an error is returned. If new telegram id is set, the telegram nickname is also to be provided and vice versa. All provided body parameters should not be blank (i.e. empty string).
//	@ID				modifyUserByJSON
//	@Tags			users
//	@Accept			json
//	@Produce		json
//
//	@Param          user body		api.EditMyProfileRequest	true	"EditMyProfileRequest"
//	@Success		200
//	@Failure		400	{object}	api.errResponse
//	@Failure		401	{object}	api.errResponse
//	@Failure		404	{object}	api.errResponse
//	@Failure		503	{object}	api.errResponse
//	@Router			/user/edit [patch]
//
// @Security Bearer
func (i *Implementation) EditMyProfile(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "api.user.EditMyProfile"

		ctx := r.Context()
		requestID := middleware.GetReqID(ctx)

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", requestID),
		)
		ctx, span := i.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
		defer span.End()

		userID := auth.UserIDFromContext(ctx)
		if userID == 0 {
			span.RecordError(api.ErrNoUserID)
			span.SetStatus(codes.Error, api.ErrNoUserID.Error())
			log.Error("no user id in context", sl.Err(api.ErrNoUserID))
			api.WriteWithError(w, http.StatusUnauthorized, api.ErrNoAuth.Error())
			return
		}

		span.AddEvent("userID extracted from context", trace.WithAttributes(attribute.Int64("id", userID)))

		req := &api.EditMyProfileRequest{}
		err := render.Bind(r, req)
		if err != nil {
			if errors.As(err, api.ValidateErr) {
				validateErr := err.(validator.ValidationErrors)
				span.RecordError(validateErr)
				span.SetStatus(codes.Error, err.Error())
				log.Error("some of the required values were not received", sl.Err(validateErr))
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
		log.Info("request body decoded", slog.Any("req", req))

		mod := convert.ToUpdateUserInfo(req, userID)

		span.AddEvent("converted to user model")

		err = i.user.EditUser(ctx, mod)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			log.Error("failed to edit user", sl.Err(err))
			api.WriteWithError(w, GetErrorCode(err), err.Error())
			return
		}

		span.AddEvent("user info updated")
		log.Info("user info updated", slog.Any("id: ", userID))

		api.WriteWithStatus(w, http.StatusOK, nil)
	}
}
