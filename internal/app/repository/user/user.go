package user

import (
	"context"
	"errors"
	"event-schedule/internal/app/model"
	"event-schedule/internal/pkg/db"
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
	ErrAlreadyExists  = errors.New("user with this nickname already exists")
	ErrQuery          = errors.New("failed to execute query")
	ErrQueryBuild     = errors.New("failed to build query")
	ErrNoRowsAffected = errors.New("no database entries affected by this operation")
	ErrPgxScan        = errors.New("failed to read database response")
	ErrNoConnection   = errors.New("could not connect to database")
	pgNoConnection    = new(*pgconn.ConnectError)
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
