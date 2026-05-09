CREATE TABLE auth_providers (
  id UUID PRIMARY KEY NOT NULL,
  user_id UUID UNIQUE NOT NULL,
  provider_name VARCHAR(255) NOT NULL DEFAULT 'local',

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_auth_provider_user
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
)
