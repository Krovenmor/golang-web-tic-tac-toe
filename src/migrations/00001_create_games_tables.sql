-- +goose Up
create table if not exists users (
    uuid UUID primary key,
    username VARCHAR(30) unique not null,
    pass VARCHAR(255)
);

create table if not exists games (
    user_uuid UUID unique references users(uuid) on delete cascade,
    field JSONB not null
);

-- +goose Down
drop table if exists games;
drop table if exists users;