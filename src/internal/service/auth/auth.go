package auth

import (
	"context"

	"github.com/google/uuid"
)

type AuthService struct {
	Repo AuthRepo
}

func NewAuthService(repo AuthRepo) *AuthService {
	return &AuthService{Repo: repo}
}

func (a *AuthService) Register(ctx context.Context, username, password string) error {
	UUID := uuid.New()
	return a.Repo.RegisterUser(ctx, UUID, username, password)
}

func (a *AuthService) LogIn(ctx context.Context, username, password string) (uuid.UUID, error) {
	return a.Repo.GetUser(ctx, username, password)
}
