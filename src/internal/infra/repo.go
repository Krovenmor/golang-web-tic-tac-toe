package infra

import (
	"WebTic-tac-toe2/internal/app"
	"WebTic-tac-toe2/internal/config"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type Repository struct {
	Pool *pgxpool.Pool
}

func NewPool(lf fx.Lifecycle, conf *config.Config) (*pgxpool.Pool, error) {
	if conf == nil {
		return nil, errors.New("Conf is nil")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, conf.DbAuth)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	lf.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

func NewRepo(pool *pgxpool.Pool) (*Repository, error) {
	return &Repository{Pool: pool}, nil
}

func (r *Repository) Save(ctx context.Context, game *app.CurrentGame) error {
	gf, err := ToRepo(game)
	if err != nil {
		return fmt.Errorf("mapping error: %w", err)
	}

	ct, err := r.Pool.Exec(ctx, saveGameQuery, gf.UUID, gf.Field)
	if err != nil {
		return fmt.Errorf("failed to save game: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("failed to save game: No uuid found")
	}

	return nil
}

func (r *Repository) SavePair(ctx context.Context, game *app.CurrentGamePair) error {
	gf, err := ToRepoPair(game)
	if err != nil {
		return fmt.Errorf("mapping error: %w", err)
	}

	ct, err := r.Pool.Exec(ctx, saveGamePairQuery,
		gf.Guuid, gf.Fuuid, gf.Suuid, gf.Field, gf.State)
	if err != nil {
		return fmt.Errorf("failed to save pairgame: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("failed to save pairgame: No game uuid found")
	}

	return nil
}

func (r *Repository) Load(ctx context.Context, uuid_ string) (*app.CurrentGame, error) {
	id, err := uuid.Parse(uuid_)
	if err != nil {
		return nil, fmt.Errorf("uuid parse error: %w", err)
	}

	var field GameFieldPG
	err = r.Pool.QueryRow(ctx, getSavedGameQuery, id).Scan(&field)
	if err != nil {
		return nil, fmt.Errorf("Load QueryRow error: %w", err)
	}

	game, err := FromRepo(id, &field)
	if err != nil {
		return nil, fmt.Errorf("Troble with FromRepo: %w", err)
	}
	return game, nil
}

func (r *Repository) LoadPair(ctx context.Context, gUUID uuid.UUID) (*app.CurrentGamePair, error) {
	var field GameFieldPG
	var fUUID, sUUID uuid.UUID
	var state int

	err := r.Pool.QueryRow(ctx, getSavedGamePairQuery, gUUID).Scan(&fUUID, &sUUID, &field, &state)
	if err != nil {
		return nil, fmt.Errorf("LoadPair QueryRow error: %w", err)
	}

	game, err := FromRepoPair(gUUID, fUUID, sUUID, &field, state)
	if err != nil {
		return nil, fmt.Errorf("Troble with FromRepo: %w", err)
	}
	return game, nil
}

func (r *Repository) GetUser(ctx context.Context, username, password string) (uuid.UUID, error) {
	var UUID uuid.UUID
	err := r.Pool.QueryRow(ctx, getUserQuery, username, password).Scan(&UUID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("GetUser QueryRow error: %w", err)
	}
	return UUID, nil
}

func (r *Repository) RegisterUser(ctx context.Context, UUID uuid.UUID, username, password string) error {
	ct, err := r.Pool.Exec(ctx, registerNewUserQuery, UUID, username, password)
	if err != nil {
		return fmt.Errorf("Trouble with register new user: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("Can't create new user")
	}
	return nil
}

func (r *Repository) GetPairGamesByState(ctx context.Context, state int) ([]uuid.UUID, error) {
	var UUIDs []uuid.UUID
	rows, err := r.Pool.Query(ctx, getPairGamesUUIDsByStateQuery, state)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UUIDs, nil
		}
		return UUIDs, fmt.Errorf("GetGamesByState QueryRow error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var gUUID uuid.UUID
		err := rows.Scan(&gUUID)
		if err != nil {
			return UUIDs, err
		}
		UUIDs = append(UUIDs, gUUID)
	}
	return UUIDs, nil
}

func (r *Repository) GetAllCompletedPairGames(ctx context.Context, uUUID uuid.UUID) ([]app.GamePairInfo, error) {
	var games []app.GamePairInfo
	rows, err := r.Pool.Query(ctx, getPairGamesInfoQuery, uUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return games, nil
		}
		return games, fmt.Errorf("getPairGamesInfoQuery QueryRow error: %w", err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[GamePairInfoPG])
	if err != nil {
		return games, err
	}
	games = make([]app.GamePairInfo, len(users))
	for idx := range users {
		games[idx] = *ToGamePairInfo(&users[idx])
	}
	return games, nil
}

func (r *Repository) GetLeaderBoard(ctx context.Context, limit uint) ([]app.GameLeaderBoardEntry, error) {
	var games []app.GameLeaderBoardEntry
	rows, err := r.Pool.Query(ctx, getLeaderBoardQuery, limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return games, nil
		}
		return games, fmt.Errorf("getLeaderBoardQuery QueryRow error: %w", err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[GameLeaderBoardPG])
	if err != nil {
		return games, err
	}
	games = make([]app.GameLeaderBoardEntry, len(users))
	for idx := range users {
		games[idx] = *ToGameLeaderBoardEntry(&users[idx])
	}
	return games, nil
}

func (r *Repository) JoinPairGame(ctx context.Context, gUUID, sUUID uuid.UUID) error {
	ct, err := r.Pool.Exec(ctx, joinGamePairQuery, gUUID, sUUID)
	if err != nil {
		return fmt.Errorf("Trouble with join new user: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("Can't join user")
	}
	return nil
}

func (r *Repository) IsPairGame(ctx context.Context, gUUID uuid.UUID) (bool, error) {
	var val int
	err := r.Pool.QueryRow(ctx, isPairGameQuery, gUUID).Scan(&val)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("Trouble with isPairGameQuery: %w", err)
	}
	if val != 1 {
		return false, nil
	}
	return true, nil
}

func (r *Repository) GetPlayerName(ctx context.Context, pUUID uuid.UUID) (string, error) {
	var val string
	err := r.Pool.QueryRow(ctx, getPlayerNameQuery, pUUID).Scan(&val)
	if err != nil {
		return "", fmt.Errorf("Trouble with getPlayerNameQuery: %w", err)
	}
	return val, nil
}

func (r *Repository) SaveRefreshToken(ctx context.Context, uUUID uuid.UUID, rToken uuid.UUID, expAt time.Time) error {
	ct, err := r.Pool.Exec(ctx, saveRefreshTokenQuery, uUUID, rToken, expAt)
	if err != nil {
		return fmt.Errorf("Trouble with saving token: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("Can't save token")
	}
	return nil
}

func (r *Repository) LoadRefreshToken(ctx context.Context, uUUID uuid.UUID, rToken uuid.UUID) (time.Time, error) {
	var expAt time.Time
	err := r.Pool.QueryRow(ctx, loadRefreshTokenQuery, uUUID, rToken).Scan(&expAt)
	if err != nil {
		return expAt, fmt.Errorf("Trouble with loadRefreshTokenQuery: %w", err)
	}
	return expAt, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, uUUID uuid.UUID, rToken uuid.UUID) error {
	ct, err := r.Pool.Exec(ctx, deleteRefreshTokenQuery, uUUID, rToken)
	if err != nil {
		return fmt.Errorf("Trouble with deleting token: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("Can't delete token")
	}
	return nil
}

func (r *Repository) DeleteExpiredRefreshTokens(ctx context.Context, uUUID uuid.UUID) error {
	_, err := r.Pool.Exec(ctx, deleteExpiredRefreshTokensQuery, uUUID)
	if err != nil {
		return fmt.Errorf("Trouble with deleting expired tokens: %w", err)
	}
	return nil
}
