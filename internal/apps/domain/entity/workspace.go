package entity

import "time"

type Workspace struct {
	ID          string
	OwnerUserID string
	Name        string
	Slug        string
	AvatarURL   string

	CreatedAt time.Time
	UpdatedAt time.Time
}
