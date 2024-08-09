-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "admin" (
  id bigserial primary key,
  username varchar(64) unique not null,
  password varchar(256) not null,
  is_active bool not null default true,
  created_at timestamptz not null default now(),
  deleted_at timestamptz
);

CREATE INDEX "admin_username" on "admin" (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "admin";
-- +goose StatementEnd
