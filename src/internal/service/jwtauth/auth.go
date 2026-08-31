package jwtauth

import (
	"WebTic-tac-toe2/internal/config"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtProvider struct {
	prvKey []byte
	repo   JwtRepoService
	aTTL   time.Duration
	rTTL   time.Duration
}

func NewJwtProvider(conf *config.Config, repo JwtRepoService) *JwtProvider {
	return &JwtProvider{
		repo:   repo,
		prvKey: conf.JWTkey,
		aTTL:   conf.JWTaccessTTL,
		rTTL:   conf.JWTrefreshTTL,
	}
}

func (p *JwtProvider) GenAccessToken(uUUID uuid.UUID) (string, error) {
	mc := jwt.MapClaims{
		"usr": uUUID.String(),
		"exp": time.Now().Add(p.aTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mc)

	return token.SignedString(p.prvKey)
}

func splitRToken(rToken string) (uuid.UUID, uuid.UUID, error) {
	spl := strings.Split(rToken, ":")
	if len(spl) != 2 {
		return uuid.Nil, uuid.Nil, errors.New("Not right token format")
	}

	sUUID, sToken := spl[0], spl[1]

	uUUID, err := uuid.Parse(sUUID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("Invalid user UUID")
	}

	tUUID, err := uuid.Parse(sToken)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("Invalid token UUID")
	}

	return uUUID, tUUID, nil
}

func (p *JwtProvider) GenRefreshToken(ctx context.Context, uUUID uuid.UUID, oldtoken string) (string, error) {
	if oldtoken != "" {
		uUUID, tUUID, err := splitRToken(oldtoken)
		if err != nil {
			return "", err
		}
		err = p.repo.DeleteRefreshToken(ctx, uUUID, tUUID)
		if err != nil {
			return "", err
		}
	}
	err := p.repo.DeleteExpiredRefreshTokens(ctx, uUUID)
	if err != nil {
		return "", err
	}

	token := uuid.New()
	expAt := time.Now().Add(p.rTTL).UTC()
	err = p.repo.SaveRefreshToken(ctx, uUUID, token, expAt)
	if err != nil {
		return "", err
	}
	return uUUID.String() + ":" + token.String(), nil
}

func (p *JwtProvider) IsValidAccess(aToken string) error {
	token, err := jwt.Parse(aToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method, got: %v", t.Header["alg"])
		}
		return p.prvKey, nil
	})

	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("Not valid token")
	}

	return nil
}

func (p *JwtProvider) IsValidRefresh(ctx context.Context, rToken string) error {
	uUUID, tUUID, err := splitRToken(rToken)
	if err != nil {
		return err
	}

	expAt, err := p.repo.LoadRefreshToken(ctx, uUUID, tUUID)
	if err != nil {
		return errors.New("Session not found")
	}

	log.Printf("[DEBUG] time.Now(): %v, expAt: %v\n", time.Now().UTC(), expAt)
	if time.Now().UTC().After(expAt) {
		log.Printf("Token expired!\n")
		return errors.New("Token expired")
	}
	log.Printf("Token not expired!\n")

	return nil
}

func (p *JwtProvider) GetUUIDfromAccessToken(aToken string) (uuid.UUID, error) {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, err := parser.Parse(aToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method, got: %v", t.Header["alg"])
		}
		return p.prvKey, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid token claims")
	}

	usr, ok := claims["usr"].(string)
	if !ok || usr == "" {
		return uuid.Nil, errors.New("UUID not found in token")
	}

	uUUID, err := uuid.Parse(usr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user UUID: %w", err)
	}

	return uUUID, nil
}

func (p *JwtProvider) GetUUIDfromRefreshToken(rToken string) (uuid.UUID, error) {
	spl := strings.Split(rToken, ":")
	if len(spl) != 2 {
		return uuid.Nil, errors.New("Not right token format")
	}

	sUUID := spl[0]

	uUUID, err := uuid.Parse(sUUID)
	if err != nil {
		return uuid.Nil, errors.New("Invalid user UUID")
	}

	return uUUID, nil
}
