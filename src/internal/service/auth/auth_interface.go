package auth

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, username, password string) error
	LogIn(ctx context.Context, username, password string) (uuid.UUID, error)
}
