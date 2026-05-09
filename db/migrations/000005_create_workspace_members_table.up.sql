CREATE TABLE workspace_members (
    id UUID PRIMARY KEY NOT NULL,

    workspace_id UUID NOT NULL,
    user_id UUID NOT NULL,

    role VARCHAR(50) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_workspace_members_workspace
        FOREIGN KEY (workspace_id)
        REFERENCES workspaces(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_workspace_members_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_workspace_member
        UNIQUE(workspace_id, user_id)
);