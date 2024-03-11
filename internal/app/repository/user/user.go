package user

import (
	"booking-schedule/internal/app/model"
	"booking-schedule/internal/pkg/db"
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	GetUserByNickname(ctx context.Context, nickName string) (*model.User, error)
	DeleteUser(ctx context.Context, userID int64) error
}

var (
	ErrNotFound       = errors.New("no user with this id")
	ErrAlreadyExists  = errors.New("this user already exists")
	ErrQuery          = errors.New("failed to execute query")
	ErrQueryBuild     = errors.New("failed to build query")
	ErrNoRowsAffected = errors.New("no database entries affected by this operation")
	ErrPgxScan        = errors.New("failed to read database response")
	ErrDuplicate      = &pgconn.PgError{
		Severity:       "ERROR",
		Code:           "23505",
		Message:        "duplicate key value violates unique constraint",
		ConstraintName: "users_telegram_id_key",
	}
	ErrNoConnection = errors.New("could not connect to database")
	pgNoConnection  = new(*pgconn.ConnectError)
)

type repository struct {
	client db.Client
	log    *slog.Logger
}

func NewUserRepository(client db.Client, log *slog.Logger) Repository {
	return &repository{
		client: client,
		log:    log,
	}
}
