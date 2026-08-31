package infra

import (
	"WebTic-tac-toe2/internal/app"
	"errors"

	"github.com/google/uuid"
)

func ToRepo(g *app.CurrentGame) (*GamePG, error) {
	if g == nil {
		return nil, errors.New("ToRepo(): game is nil")
	}
	ng := &GamePG{}
	ng.Field.Field = g.GField.Field
	ng.UUID = g.UUID
	return ng, nil
}

func ToRepoPair(g *app.CurrentGamePair) (*GamePairPG, error) {
	if g == nil {
		return nil, errors.New("ToRepoPair(): game is nil")
	}
	return &GamePairPG{
		Field: GameFieldPG{Field: g.GField.Field},
		Guuid: g.Guuid,
		Fuuid: g.Fuuid,
		Suuid: g.Suuid,
		State: g.State,
	}, nil
}

func FromRepo(UUID uuid.UUID, field *GameFieldPG) (*app.CurrentGame, error) {
	if field == nil {
		return nil, errors.New("FromRepo(): field is nil")
	}
	ng := &app.CurrentGame{}
	ng.GField.Field = field.Field
	ng.UUID = UUID
	return ng, nil
}

func FromRepoPair(gUUID, fUUID, sUUID uuid.UUID, field *GameFieldPG, state int) (*app.CurrentGamePair, error) {
	if field == nil {
		return nil, errors.New("FromRepoPair(): field is nil")
	}
	return &app.CurrentGamePair{
		GField: app.GameField{Field: field.Field},
		Guuid:  gUUID,
		Fuuid:  fUUID,
		Suuid:  sUUID,
		State:  state,
	}, nil
}

func ToGamePairInfo(from *GamePairInfoPG) *app.GamePairInfo {
	return &app.GamePairInfo{
		Guuid:   from.Guuid,
		Wuuid:   from.Wuuid,
		State:   from.State,
		Created: from.Created,
	}
}

func ToGameLeaderBoardEntry(from *GameLeaderBoardPG) *app.GameLeaderBoardEntry {
	return &app.GameLeaderBoardEntry{
		PlayerUUID:  from.PlayerUUID,
		PlayerLogin: from.PlayerLogin,
		Ratio:       from.Ratio,
	}
}
