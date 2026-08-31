-- +goose Up
create table if not exists UserTokens (
    uUUID uuid references users(uuid) on delete cascade,
    tUUID uuid not null,
    expAt timestamp not null,

    primary key (uUUID, tUUID)
);

-- +goose Down
drop table if exists UserTokens;