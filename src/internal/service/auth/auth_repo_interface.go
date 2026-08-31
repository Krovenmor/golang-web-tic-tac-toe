package auth

import (
	"context"

	"github.com/google/uuid"
)

type AuthRepo interface {
	GetUser(ctx context.Context, username, password string) (uuid.UUID, error)
	RegisterUser(ctx context.Context, UUID uuid.UUID, username, password string) error
}
