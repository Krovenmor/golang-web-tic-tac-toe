package gamehandler

import (
	"WebTic-tac-toe2/internal/app"
	"errors"

	"github.com/google/uuid"
)

func FromWeb(gs *GameCreateResponseBody) (*app.CurrentGame, error) {
	if gs == nil {
		return nil, errors.New("FromWeb(): gs is nil")
	}
	UUID, err := uuid.Parse(gs.GUUID)
	if err != nil {
		return nil, err
	}
	return &app.CurrentGame{UUID: UUID, GField: app.GameField(gs.GameField)}, nil
}

func FromWebPair(gs *GameCreateResponseBody, fUUID, sUUID uuid.UUID, state int) (*app.CurrentGamePair, error) {
	if gs == nil {
		return nil, errors.New("FromWeb(): gs is nil")
	}
	UUID, err := uuid.Parse(gs.GUUID)
	if err != nil {
		return nil, err
	}
	return &app.CurrentGamePair{
		Guuid:  UUID,
		Fuuid:  fUUID,
		Suuid:  sUUID,
		GField: app.GameField(gs.GameField),
		State:  state,
	}, nil
}

func ToWeb(cg *app.CurrentGame) (*GameCreateResponseBody, error) {
	if cg == nil {
		return nil, errors.New("ToWeb(): cg is nil")
	}
	return &GameCreateResponseBody{GUUID: cg.UUID.String(),
		GameField: GameFieldWeb{cg.GField.Field},
		State:     cg.State}, nil
}

func ToWebPair(cg *app.CurrentGamePair) (*GamePairCreateResponseBody, error) {
	if cg == nil {
		return nil, errors.New("ToWeb(): cg is nil")
	}
	return &GamePairCreateResponseBody{
		GUUID:     cg.Guuid,
		FUUID:     cg.Fuuid,
		SUUID:     cg.Suuid,
		GameField: GameFieldWeb{cg.GField.Field},
		State:     cg.State}, nil
}

func ToWebPairList(games []uuid.UUID) (*GamePairListResponseBody, error) {
	return &GamePairListResponseBody{
		Games: games,
	}, nil
}

func ToWebPairCompletedList(games []app.GamePairInfo) ([]GamePairInfoResponseBody, error) {
	l := make([]GamePairInfoResponseBody, len(games))
	for idx, val := range games {
		l[idx].GUUID = val.Guuid
		l[idx].WUUID = val.Wuuid
		l[idx].State = val.State
		l[idx].Created = val.Created
	}
	return l, nil
}

func ToWebLeaderBoard(lb []app.GameLeaderBoardEntry) ([]GameLeaderBoardEntryResponseBody, error) {
	l := make([]GameLeaderBoardEntryResponseBody, len(lb))
	for idx, val := range lb {
		l[idx].PlayerUUID = val.PlayerUUID
		l[idx].Ratio = val.Ratio
		l[idx].PlayerLogin = val.PlayerLogin
	}
	return l, nil
}

func ToWebShortSolo(cg *app.CurrentGame) (*GameShortInfoResponseBody, error) {
	if cg == nil {
		return nil, errors.New("ToWebShort(): cg is nil")
	}
	return &GameShortInfoResponseBody{
		GameField: GameFieldWeb{Field: cg.GField.Field},
		State:     cg.State,
	}, nil
}

func ToWebShortPair(cg *app.CurrentGamePair) (*GameShortInfoResponseBody, error) {
	if cg == nil {
		return nil, errors.New("ToWebShort(): cg is nil")
	}
	return &GameShortInfoResponseBody{
		GameField: GameFieldWeb{Field: cg.GField.Field},
		State:     cg.State,
	}, nil
}

func FromWebShortSolo(field GameFieldWeb, from *app.CurrentGame) (*app.CurrentGame, error) {
	if from == nil {
		return nil, errors.New("FromWebShortSolo(): from is nil")
	}
	return &app.CurrentGame{
		GField: app.GameField{Field: field.Field},
		UUID:   from.UUID,
		State:  from.State,
	}, nil
}

func FromWebShortPair(field GameFieldWeb, from *app.CurrentGamePair) (*app.CurrentGamePair, error) {
	if from == nil {
		return nil, errors.New("FromWebShortPair(): from is nil")
	}
	return &app.CurrentGamePair{
		GField: app.GameField{Field: field.Field},
		Guuid:  from.Guuid,
		Fuuid:  from.Fuuid,
		Suuid:  from.Suuid,
		State:  from.State,
	}, nil
}

func ToPlayerInfo(playerName string) *PlayerInfoResponseBody {
	return &PlayerInfoResponseBody{
		PlayerName: playerName,
	}
}
