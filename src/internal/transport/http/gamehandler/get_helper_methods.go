package gamehandler

import (
	"WebTic-tac-toe2/internal/app"
	auth "WebTic-tac-toe2/internal/transport/http/authhandler"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func getUuidFromPath(req *http.Request) (uuid.UUID, error) {
	UUID := req.PathValue("uuid")
	if UUID == "" {
		return uuid.Nil, errors.New("Empty gUUID")
	}
	gUUID, err := uuid.Parse(UUID)
	if err != nil {
		return uuid.Nil, errors.New("Bad gUUID")
	}
	return gUUID, nil
}

func getUuidFromBody(req *http.Request) (uuid.UUID, error) {
	var ng UuidIncomeBody
	err := json.NewDecoder(req.Body).Decode(&ng)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Bad UUID: %q", err.Error())
	}
	return ng.UUID, nil
}

func getUuidFromAuth(req *http.Request) (uuid.UUID, error) {
	val, ok := req.Context().Value(auth.UserUUIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("Trouble with receiving uUUID from UserAuthenticator")
	}
	return val, nil
}

func (gh *GameHandler) getPairGame(ctx context.Context, gUUID uuid.UUID) (*app.CurrentGamePair, error) {
	game, err := gh.service.GetPairGame(ctx, gUUID)

	if err != nil {
		return nil, err
	}
	if game == nil {
		return nil, errors.New("getPairGame: game == nil")
	}

	return game, nil
}

func (gh *GameHandler) getSoloGame(ctx context.Context, gUUID uuid.UUID) (*app.CurrentGame, error) {
	game, err := gh.service.GetSoloGame(ctx, gUUID)

	if err != nil {
		return nil, err
	}
	if game == nil {
		return nil, errors.New("getSoloGame: game == nil")
	}

	return game, nil
}
