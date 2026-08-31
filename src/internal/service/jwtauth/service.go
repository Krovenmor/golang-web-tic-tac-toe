package jwtauth

import (
	"context"

	"github.com/google/uuid"
)

type JWTService interface {
	GenAccessToken(uUUID uuid.UUID) (string, error)
	GenRefreshToken(ctx context.Context, uUUID uuid.UUID, oldtoken string) (string, error)

	IsValidAccess(token string) error
	IsValidRefresh(ctx context.Context, token string) error

	GetUUIDfromAccessToken(token string) (uuid.UUID, error)
	GetUUIDfromRefreshToken(token string) (uuid.UUID, error)
}
