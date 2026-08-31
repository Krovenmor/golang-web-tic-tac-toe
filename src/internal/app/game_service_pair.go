package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

const (
	WAIT_STATE = iota
	FMOVING_STATE
	SMOVING_STATE
	DRAW_STATE
	FWIN_STATE
	SWIN_STATE
)

func (gs *gameService) NewPairGame(ctx context.Context, fUUID uuid.UUID) (*CurrentGamePair, error) {
	gUUID := uuid.New()
	ng := &CurrentGamePair{Guuid: gUUID, Fuuid: fUUID, Suuid: uuid.Nil, GField: GameField{}, State: WAIT_STATE}
	err := gs.repo.SavePair(ctx, ng)
	if err != nil {
		return nil, err
	}
	return ng, nil
}

func (gs *gameService) MakeMovePair(ctx context.Context, game *CurrentGamePair) (*CurrentGamePair, error) {
	if game == nil {
		return nil, errors.New("MakeMovePair(): Game is nil")
	}

	switch game.State {
	case WAIT_STATE, SMOVING_STATE:
		game.State = FMOVING_STATE
	case FMOVING_STATE:
		game.State = SMOVING_STATE
	}

	isfull := isFull(&game.GField)
	if isfull {
		game.State = DRAW_STATE
	}

	who, iswin := isWin(&game.GField)
	if iswin {
		if who == CELL_PLAYER {
			game.State = FWIN_STATE
		} else {
			game.State = SWIN_STATE
		}
	}

	err := gs.repo.SavePair(ctx, game)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (gs *gameService) IsValidPair(ctx context.Context, game *CurrentGamePair) (bool, error) {
	if game == nil {
		return false, errors.New("IsValidPair(): game == nil")
	}

	g, err := gs.repo.LoadPair(ctx, game.Guuid)
	if err != nil {
		return false, err
	}
	if g == nil {
		return false, errors.New("IsValid(): Loaded Game is nil")
	}
	if g.State == WAIT_STATE {
		return true, nil
	}
	diff_count := 0

	player_cell := CELL_EMPTY
	switch g.State {
	case FMOVING_STATE:
		player_cell = CELL_PLAYER
	case SMOVING_STATE:
		player_cell = CELL_COMPUTER
	}

	for i, sl := range game.GField.Field {
		for j, val := range sl {
			if val < CELL_EMPTY || val > CELL_COMPUTER {
				return false, errors.New("Not valid val on field")
			}
			f_val := g.GField.Field[i][j]
			if val != f_val {
				if f_val != CELL_EMPTY {
					return false, errors.New("Not empty cell can't be changed")
				}
				if val != player_cell {
					return false, errors.New("Not right player moved")
				}
				diff_count++
			}
			if diff_count > 1 {
				return false, errors.New("Difference is too big")
			}
		}
	}
	return true, nil
}

func (gs *gameService) IsOverPair(ctx context.Context, game *CurrentGamePair) (bool, error) {
	if game == nil {
		return false, errors.New("IsOverPair(): game == nil")
	}

	switch game.State {
	case FWIN_STATE, SWIN_STATE, DRAW_STATE:
		return true, nil
	default:
		return false, nil
	}
}

func (gs *gameService) GetAllAvailablePairGames(ctx context.Context) ([]uuid.UUID, error) {
	games, err := gs.repo.GetPairGamesByState(ctx, WAIT_STATE)
	if err != nil {
		return games, err
	}
	return games, nil
}

func (gs *gameService) GetAllCompletedPairGames(ctx context.Context, uUUID uuid.UUID) ([]GamePairInfo, error) {
	return gs.repo.GetAllCompletedPairGames(ctx, uUUID)
}

func (gs *gameService) GetLeaderBoard(ctx context.Context, limit uint) ([]GameLeaderBoardEntry, error) {
	return gs.repo.GetLeaderBoard(ctx, limit)
}

func (gs *gameService) JoinPairGame(ctx context.Context, gUUID, sUUID uuid.UUID) (*CurrentGamePair, error) {
	err := gs.repo.JoinPairGame(ctx, gUUID, sUUID)
	if err != nil {
		return nil, err
	}
	game, err := gs.repo.LoadPair(ctx, gUUID)
	if err != nil {
		return nil, err
	}
	game.State = FMOVING_STATE
	err = gs.repo.SavePair(ctx, game)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (gs *gameService) IsPairGame(ctx context.Context, gUUID uuid.UUID) (bool, error) {
	return gs.repo.IsPairGame(ctx, gUUID)
}

func (gs *gameService) GetPairGame(ctx context.Context, gUUID uuid.UUID) (*CurrentGamePair, error) {
	return gs.repo.LoadPair(ctx, gUUID)
}
