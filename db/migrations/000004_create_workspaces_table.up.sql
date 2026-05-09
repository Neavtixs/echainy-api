CREATE TABLE workspaces (
  id UUID PRIMARY KEY NOT NULL,
  owner_user_id UUID NOT NULL,
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL,
  avatar_url VARCHAR(255),

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_workspace_user
    FOREIGN KEY (owner_user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,

  CONSTRAINT unique_workspace_owner_slug
    UNIQUE (owner_user_id, slug)
);
