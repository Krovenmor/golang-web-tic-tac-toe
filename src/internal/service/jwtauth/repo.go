package jwtauth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type JwtRepoService interface {
	SaveRefreshToken(ctx context.Context, uUUID uuid.UUID, rToken uuid.UUID, expAt time.Time) error
	LoadRefreshToken(ctx context.Context, uUUID uuid.UUID, rToken uuid.UUID) (time.Time, error)
	DeleteRefreshToken(ctx context.Context, uUUID uuid.UUID, rToken uuid.UUID) error
	DeleteExpiredRefreshTokens(ctx context.Context, uUUID uuid.UUID) error
}
