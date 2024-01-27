package schedule

import (
	"context"

	gofakeit "github.com/brianvoe/gofakeit/v6"
)

func (r *repository) AddEvent(ctx context.Context) (string, error) {
	return gofakeit.UUID(), nil
}
