-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
  id bigserial not null,
  username varchar(256) unique not null,
  external_id varchar(256),
  user_type varchar(64) not null,
  referral_code varchar(32) not null unique,
  referral varchar(32) references "user"(referral_code),
  banned_at timestamptz,
  created_at timestamptz not null default now(),
  deleted_at timestamptz
);

CREATE INDEX "user_username" on "user" (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
