package app

import (
	"context"

	"github.com/google/uuid"
)

type GameService interface {
	NewGame(ctx context.Context, UUID uuid.UUID) (*CurrentGame, error)
	MakeMove(ctx context.Context, game *CurrentGame) (*CurrentGame, error)
	IsValid(ctx context.Context, game *CurrentGame) (bool, error)
	IsOver(ctx context.Context, game *CurrentGame) (bool, error)

	NewPairGame(ctx context.Context, fUUID uuid.UUID) (*CurrentGamePair, error)
	MakeMovePair(ctx context.Context, game *CurrentGamePair) (*CurrentGamePair, error)
	IsValidPair(ctx context.Context, game *CurrentGamePair) (bool, error)
	IsOverPair(ctx context.Context, game *CurrentGamePair) (bool, error)

	JoinPairGame(ctx context.Context, gUUID, sUUID uuid.UUID) (*CurrentGamePair, error)
	IsPairGame(ctx context.Context, gUUID uuid.UUID) (bool, error)

	GetSoloGame(ctx context.Context, gUUID uuid.UUID) (*CurrentGame, error)
	GetPairGame(ctx context.Context, gUUID uuid.UUID) (*CurrentGamePair, error)
	GetPlayerName(ctx context.Context, pUUID uuid.UUID) (string, error)

	GetAllAvailablePairGames(ctx context.Context) ([]uuid.UUID, error)
	GetAllCompletedPairGames(ctx context.Context, uUUID uuid.UUID) ([]GamePairInfo, error)
	GetLeaderBoard(ctx context.Context, limit uint) ([]GameLeaderBoardEntry, error)
}
