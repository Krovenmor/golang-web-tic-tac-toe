package infra

// Game queries
const (
	saveGameQuery = `
		INSERT INTO games (user_uuid, field) 
		SELECT users.uuid, $2
		FROM users
		WHERE users.uuid = $1
		ON CONFLICT (user_uuid) 
		DO UPDATE SET field = EXCLUDED.field;
	`

	saveGamePairQuery = `
		INSERT INTO pair_games (uuid, user_uuid_f, user_uuid_s, field, state)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (uuid)
		DO UPDATE SET field = EXCLUDED.field, state = EXCLUDED.state, user_uuid_s = EXCLUDED.user_uuid_s;
	`

	joinGamePairQuery = `
		update pair_games
		set user_uuid_s = $2
		where uuid = $1 and user_uuid_f != $2;
	`

	isPairGameQuery = `
		select 1
		from pair_games
		where uuid = $1
	`

	getPairGamesUUIDsByStateQuery = `
		select uuid
		from pair_games
		where state = $1
	`

	getPairGamesInfoQuery = `
		select
			uuid,
			case when state = 4 then user_uuid_f
				when state = 5 then user_uuid_s
			end as winner_uuid,
			state,
			created
		from pair_games
		where state > 2 and $1 in (user_uuid_f, user_uuid_s)
		order by created desc;
	`

	getLeaderBoardQuery = `
		with tmp as (
			select
				user_uuid_f as player,
				u.username as login,
				(case when state = 4 then 1.0 else 0.0 end) as is_win
			from pair_games p
			join users u on u.uuid = p.user_uuid_f
			where state in (3, 4, 5)
			union all
			select
				user_uuid_s as player,
				u.username as login,
				(case when state = 5 then 1.0 else 0.0 end) as is_win
			from pair_games p
			join users u on u.uuid = p.user_uuid_s
			where state in (3, 4, 5)
		)
		select
			player,
			login,
			coalesce(sum(is_win) / nullif(count(*) - sum(is_win), 0), sum(is_win)) as ratio
		from tmp
		group by player, login
		order by ratio DESC
		Limit $1;
	`

	getSavedGameQuery = `
		select field 
		from games
		where user_uuid = $1;
	`

	getSavedGamePairQuery = `
		select user_uuid_f, user_uuid_s, field, state
		from pair_games
		where uuid = $1;
	`
)

// Auth queries
const (
	registerNewUserQuery = `
		insert into users (uuid, username, pass)
		values ($1, $2, $3);
	`

	getUserQuery = `
		select uuid
		from users
		where username = $1 and pass = $2;
	`

	getPlayerNameQuery = `
		select username
		from users
		where uuid = $1;
	`
)

// JWT queries
const (
	saveRefreshTokenQuery = `
		insert into UserTokens (uUUID, tUUID, expAt)
		values ($1, $2, $3);
	`

	loadRefreshTokenQuery = `
		select expAt
		from UserTokens
		where uUUID = $1 and tUUID = $2;
	`

	deleteRefreshTokenQuery = `
		delete 
		from UserTokens
		where uUUID = $1 and tUUID = $2;
	`

	deleteExpiredRefreshTokensQuery = `
		delete 
		from UserTokens
		where uUUID = $1 and expAt < Now();
	`
)
