package jwt

import (
	"booking-schedule/internal/logger/sl"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	ErrNoID            = errors.New("user id not set")
	ErrInvalidToken    = errors.New("invalid token")

	ErrParseID  = errors.New("parsing user id failed")
	ErrParseExp = errors.New("parsing token expiration failed")
)

// GenerateToken takes a user ID and
func (s *service) GenerateToken(ctx context.Context, userID int64) (string, error) {
	const op = "service.jwt.GenerateToken"

	requestID := middleware.GetReqID(ctx)

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)

	_, span := s.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
	defer span.End()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(s.expiration).Unix(),
	})

	span.AddEvent("token generated", trace.WithAttributes(attribute.Int64("user_id", userID)))
	log.Info("token generated", slog.Int64("user id:", userID))

	return token.SignedString([]byte(s.jwtSecret))
}

// VerifyToken parses and validates a jwt token. It returns the userID if the token is valid.
func (s *service) VerifyToken(ctx context.Context, tokenString string) (int64, error) {
	const op = "service.jwt.VerifyToken"

	requestID := middleware.GetReqID(ctx)

	log := s.log.With(
		slog.String("op", op),
		slog.String("request_id", requestID),
	)

	_, span := s.tracer.Start(ctx, op, trace.WithAttributes(attribute.String("request_id", requestID)))
	defer span.End()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			span.RecordError(ErrUnsupportedSign)
			span.SetStatus(codes.Error, ErrUnsupportedSign.Error())
			log.Error("unexpected signing method: ", token.Header["alg"])
			return nil, ErrUnsupportedSign
		}
		return []byte(s.jwtSecret), nil
	}, jwt.WithJSONNumber())

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("parsing token failed: ", sl.Err(err))
		return 0, err
	}

	span.AddEvent("token parsed")

	claims, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		span.RecordError(ErrInvalidToken)
		span.SetStatus(codes.Error, ErrInvalidToken.Error())
		log.Error("invalid token", sl.Err(ErrInvalidToken))
		return 0, ErrInvalidToken
	}

	userID := claims["userID"]
	userIDInt, err := userID.(json.Number).Int64()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("issue parsing user id", sl.Err(err))
		return 0, ErrParseID

	}

	if userIDInt == 0 {
		span.RecordError(ErrNoID)
		span.SetStatus(codes.Error, ErrNoID.Error())
		log.Error("empty user id", sl.Err(ErrNoID))
		return 0, ErrNoID
	}

	span.AddEvent("userID acquired")

	exp, err := claims["exp"].(json.Number).Int64()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error("issue parsing token expiration", sl.Err(err))
		return 0, ErrParseExp

	}

	span.AddEvent("token expiration acquired")

	if exp < time.Now().Unix() {
		span.RecordError(jwt.ErrTokenExpired)
		span.SetStatus(codes.Error, jwt.ErrTokenExpired.Error())
		log.Error("token expired", sl.Err(jwt.ErrTokenExpired))
		return 0, jwt.ErrTokenExpired
	}

	return userIDInt, nil

}
