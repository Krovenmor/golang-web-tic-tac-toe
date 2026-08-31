package app

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

const (
	DEF_STATUS = iota
	COMP_WIN_STATUS
	PLYR_WIN_STATUS
	FULL_FLD_STATUS
)

type gameService struct {
	repo GameRepository
}

func NewGameService(repo GameRepository) GameService {
	return &gameService{repo: repo}
}

func (gs *gameService) NewGame(ctx context.Context, UUID uuid.UUID) (*CurrentGame, error) {
	ng := &CurrentGame{UUID: UUID, GField: GameField{}, State: DEF_STATUS}
	err := gs.repo.Save(ctx, ng)
	if err != nil {
		return nil, err
	}
	return ng, nil
}

func score(field *GameField) int {
	winner, is_win := isWin(field)
	if is_win {
		if winner == CELL_COMPUTER {
			return 10
		}
		return -10
	}
	return 0
}

func isOver(field *GameField) (int, bool) {
	sc := score(field)
	if sc != 0 {
		return sc, true
	}
	return sc, isFull(field)
}

type Vector2i struct {
	X, Y int
}

func calc(ctx context.Context, field GameField, isCompTurn bool) (int, error) {
	select {
	case <-ctx.Done():
		return -1, ctx.Err()
	default:
		sc, isOver := isOver(&field)
		if isOver {
			return sc, nil
		}

		if isCompTurn {
			maxVal := -20
			l := len(field.Field)
			for i := range l {
				for j := range l {
					if field.Field[i][j] == CELL_EMPTY {
						field.Field[i][j] = CELL_COMPUTER

						scr, err := calc(ctx, field, !isCompTurn)
						if err != nil {
							return -1, err
						}
						if scr > maxVal {
							maxVal = scr
						}

						field.Field[i][j] = CELL_EMPTY
					}
				}
			}
			return maxVal, nil
		} else {
			minVal := 20
			l := len(field.Field)
			for i := range l {
				for j := range l {
					if field.Field[i][j] == CELL_EMPTY {
						field.Field[i][j] = CELL_PLAYER

						scr, err := calc(ctx, field, !isCompTurn)
						if err != nil {
							return -1, err
						}
						if scr < minVal {
							minVal = scr
						}

						field.Field[i][j] = CELL_EMPTY
					}
				}
			}
			return minVal, nil
		}
	}
}

func (gs *gameService) makeNextMove(ctx context.Context, game *CurrentGame) (bool, error) {
	if isFull(&game.GField) {
		return false, nil
	}

	maxScr := -20
	bPos := Vector2i{-1, -1}
	for i, sl := range game.GField.Field {
		for j, val := range sl {
			if val == CELL_EMPTY {
				game.GField.Field[i][j] = CELL_COMPUTER

				scr, err := calc(ctx, game.GField, false)
				if err != nil {
					return false, err
				}
				if scr > maxScr {
					maxScr = scr
					bPos = Vector2i{X: i, Y: j}
				}

				game.GField.Field[i][j] = CELL_EMPTY
			}
		}
	}
	if maxScr > -20 {
		game.GField.Field[bPos.X][bPos.Y] = CELL_COMPUTER
		return true, nil
	}
	return false, nil
}

func (gs *gameService) MakeMove(ctx context.Context, game *CurrentGame) (*CurrentGame, error) {
	if game == nil {
		return nil, errors.New("MakeMove(): Game is nil")
	}
	isMoved, err := gs.makeNextMove(ctx, game)
	if err != nil || !isMoved {
		return nil, err
	}
	err = gs.repo.Save(ctx, game)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (gs *gameService) IsValid(ctx context.Context, game *CurrentGame) (bool, error) {
	if game == nil {
		return false, errors.New("IsValid(): game == nil")
	}

	g, err := gs.repo.Load(ctx, game.UUID.String())
	if err != nil {
		return false, err
	}
	if g == nil {
		return false, errors.New("IsValid(): Loaded Game is nil")
	}
	diff_count := 0
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
				diff_count++
			}
			if diff_count > 1 {
				return false, errors.New("Difference is too big")
			}
		}
	}
	return true, nil
}

func isFull(field *GameField) bool {
	for _, sl := range field.Field {
		for _, val := range sl {
			if val == CELL_EMPTY {
				return false
			}
		}
	}
	return true
}

func checkVH(field *GameField) (int, bool) {
	l := len(field.Field)
	for i := range l {
		hor_fv, ver_fv := -1, -1
		is_hor, is_ver := true, true
		for j := range l {
			hv := field.Field[j][i]
			vv := field.Field[i][j]

			if hor_fv == -1 {
				hor_fv = hv
				ver_fv = vv
				if hor_fv == CELL_EMPTY {
					is_hor = false
				}
				if ver_fv == CELL_EMPTY {
					is_ver = false
				}
				continue
			}

			if hv != hor_fv {
				is_hor = false
			}
			if vv != ver_fv {
				is_ver = false
			}
		}

		if is_hor {
			return hor_fv, true
		}
		if is_ver {
			return ver_fv, true
		}
	}

	return -1, false
}

func checkDi(field *GameField) (int, bool) {
	if field.Field[0][0] != CELL_EMPTY && field.Field[0][0] == field.Field[1][1] && field.Field[1][1] == field.Field[2][2] {
		return field.Field[0][0], true
	}
	if field.Field[0][2] != CELL_EMPTY && field.Field[0][2] == field.Field[1][1] && field.Field[1][1] == field.Field[2][0] {
		return field.Field[0][2], true
	}
	return -1, false
}

func isWin(field *GameField) (int, bool) {
	winner, is_win := checkDi(field)
	if !is_win {
		winner, is_win = checkVH(field)
	}
	return winner, is_win
}

func (*gameService) IsOver(ctx context.Context, game *CurrentGame) (bool, error) {
	if game == nil {
		return false, errors.New("IsOver(): game == nil")
	}
	winner, is_win := isWin(&game.GField)
	is_full := isFull(&game.GField)
	if is_win || is_full {
		if is_win {
			if winner == CELL_COMPUTER {
				game.State = COMP_WIN_STATUS
			} else {
				game.State = PLYR_WIN_STATUS
			}
		} else if is_full {
			game.State = FULL_FLD_STATUS
		}
	} else {
		game.State = DEF_STATUS
	}
	return is_win || is_full, nil
}

func (gs *gameService) GetSoloGame(ctx context.Context, gUUID uuid.UUID) (*CurrentGame, error) {
	return gs.repo.Load(ctx, gUUID.String())
}

func (gs *gameService) GetPlayerName(ctx context.Context, pUUID uuid.UUID) (string, error) {
	return gs.repo.GetPlayerName(ctx, pUUID)
}
