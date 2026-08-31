package app

import (
	"context"

	"github.com/google/uuid"
)

type GameRepository interface {
	Save(ctx context.Context, game *CurrentGame) error
	Load(ctx context.Context, uuid string) (*CurrentGame, error)

	SavePair(ctx context.Context, game *CurrentGamePair) error
	LoadPair(ctx context.Context, gUUID uuid.UUID) (*CurrentGamePair, error)
	JoinPairGame(ctx context.Context, gUUID, sUUID uuid.UUID) error
	IsPairGame(ctx context.Context, gUUID uuid.UUID) (bool, error)

	GetPairGamesByState(ctx context.Context, state int) ([]uuid.UUID, error)
	GetPlayerName(ctx context.Context, pUUID uuid.UUID) (string, error)
	GetAllCompletedPairGames(ctx context.Context, uUUID uuid.UUID) ([]GamePairInfo, error)
	GetLeaderBoard(ctx context.Context, limit uint) ([]GameLeaderBoardEntry, error)
}
