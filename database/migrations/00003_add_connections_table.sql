-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "connection" (
  id bigserial not null,
  status varchar(16) not null default 'connected',
  username varchar(256) not null,
  external_id varchar(256) not null,
  remote_ip varchar(256) not null,
  location varchar(256) not null,
  user_agent varchar(256) not null,
  hostname varchar(256) not null,
  download_traffic_usage bigint not null,
  upload_traffic_usage bigint not null,
  connected_at timestamptz not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

CREATE INDEX "connection_username" on "connection" (username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "connection";
-- +goose StatementEnd
