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

type WorkspaceMemberWorkspace struct {
	ID          string
	OwnerUserID string
	Name        string
	Slug        string
	AvatarURL   string
	Role        string
}
