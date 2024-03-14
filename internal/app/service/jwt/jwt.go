package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/trace"
)

// Service is an interface that represents all the capabilities for the JWT service.
type Service interface {
	GenerateToken(ctx context.Context, userID int64) (string, error)
	VerifyToken(ctx context.Context, token string) (int64, error)
}

type service struct {
	jwtSecret  string
	expiration time.Duration
	log        *slog.Logger
	tracer     trace.Tracer
}

// New creates a service with a provided JWT secret string and expiration (hourly) number. It implements
// the JWT Service interface.
func NewJWTService(jwtSecret string, expiration time.Duration, log *slog.Logger, tracer trace.Tracer) Service {
	return &service{jwtSecret, expiration, log, tracer}
}

var (
	ErrUnsupportedSign = errors.New("unexpected signing method")
	ErrParseID         = errors.New("parsing user id failed")
	ErrParseExp        = errors.New("parsing token expiration failed")
	ErrNoID            = errors.New("user id not set")
	ErrInvalidToken    = errors.New("invalid token")
)

// GenerateToken takes a user ID and
func (s *service) GenerateToken(ctx context.Context, userID int64) (string, error) {
	const op = "auth.service.GenerateToken"

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(s.expiration).Unix(),
	})

	log.Info("generated token", slog.Any("user id:", userID))
	tkn, err := token.SignedString([]byte(s.jwtSecret))

	return tkn, err
}

// VerifyToken parses and validates a jwt token. It returns the userID if the token is valid.
func (s *service) VerifyToken(ctx context.Context, tokenString string) (int64, error) {
	const op = "jwt.service.VerifyToken"

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("unexpected signing method: ", token.Header["alg"])
			return nil, ErrUnsupportedSign
		}
		return []byte(s.jwtSecret), nil
	}, jwt.WithJSONNumber())

	if err != nil {
		log.Error("parsing token failed: ", err)
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		log.Error("invalid token")
		return 0, ErrInvalidToken
	}

	userID := claims["userID"]
	userIDInt, err := userID.(json.Number).Int64()
	if err != nil {
		log.Error("issue parsing user id", err)
		return 0, ErrParseID

	}

	if userIDInt == 0 {
		log.Error("empty user id")
		return 0, ErrNoID
	}

	exp, err := claims["exp"].(json.Number).Int64()
	if err != nil {
		log.Error("issue parsing token expiration", err)
		return 0, ErrParseExp

	}

	if exp < time.Now().Unix() {
		log.Error("token expired", jwt.ErrTokenExpired)
		return 0, jwt.ErrTokenExpired
	}

	return userIDInt, nil

}
