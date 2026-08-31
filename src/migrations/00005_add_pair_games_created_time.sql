-- +goose Up
alter table pair_games add column if not exists created timestamp default now();

-- +goose Down
alter table pair_games drop column if exists created;