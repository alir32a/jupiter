-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "package" (
  id bigserial not null,
  user_id bigint not null,
  traffic_limit bigint not null,
  download_traffic_usage bigint not null default 0,
  upload_traffic_usage bigint not null default 0,
  max_connections int not null,
  is_trial boolean not null default false,
  expiration_in_days int not null,
  expire_at timestamptz,
  created_at timestamptz not null default now()
);

CREATE UNIQUE INDEX trial_package_idx ON "package" (user_id) WHERE is_trial = true;
CREATE INDEX "package_user_id" on "package" (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "package";
-- +goose StatementEnd
