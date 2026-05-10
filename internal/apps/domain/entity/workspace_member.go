package entity

import "time"

type WorkspaceMember struct {
	ID          string
	WorkspaceID string
	UserID      string
	Role        string

	CreatedAt time.Time
	UpdatedAt time.Time
}
