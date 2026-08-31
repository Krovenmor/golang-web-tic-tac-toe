-- +goose Up
create table if not exists pair_games (
    uuid UUID PRIMARY KEY,
    user_uuid_f UUID references users(uuid) on delete cascade,
    user_uuid_s UUID,
    field JSONB not null,
    state int not null
);

-- +goose Down
drop table if exists pair_games;