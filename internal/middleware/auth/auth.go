package auth

import (
	"booking-schedule/internal/app/api"
	"booking-schedule/internal/app/service/jwt"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type userID int64

const (
	keyUserID userID = 0
)

var (
	errMissingToken = errors.New("missing bearer token")
	errInvalidToken = errors.New("token is invalid")
)

// Auth creates a middleware function that retrieves a bearer token and validates the token.
// The middleware sets the userID in the jwt payload into the request context. If the token is
// invalid, it will write an Unauthorized response.
func Auth(logger *slog.Logger, jwtService jwt.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "auth.service.Auth"

			ctx := r.Context()

			log := logger.With(
				slog.String("op", op),
				slog.String("request_id", middleware.GetReqID(ctx)),
			)

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
				log.Error("missing token ", errMissingToken)
				render.Status(r, http.StatusUnauthorized)
				api.WriteWithError(w, http.StatusUnauthorized, errMissingToken.Error())
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := jwtService.VerifyToken(ctx, token)
			if err != nil {
				log.Error("issue verifying jwt token", err)
				render.Status(r, http.StatusUnauthorized)
				api.WriteWithError(w, http.StatusUnauthorized, errInvalidToken.Error())
				return
			}

			r = r.WithContext(withUser(ctx, userID))
			next.ServeHTTP(w, r)
		})
	}
}

// UserIDFromContext returns a user ID from context
func UserIDFromContext(ctx context.Context) int64 {
	if userID, ok := ctx.Value(keyUserID).(int64); ok {
		return userID
	}

	return 0
}

// withUser adds the userID to a context object and returns that context
func withUser(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, keyUserID, userID)
}
