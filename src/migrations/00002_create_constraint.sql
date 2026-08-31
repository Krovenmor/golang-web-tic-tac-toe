-- +goose Up
alter table games add constraint user_uuid_unique unique (user_uuid);

-- +goose Down
alter table games drop constraint user_uuid_unique;