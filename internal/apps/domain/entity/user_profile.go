package entity

import "time"

type UserProfile struct {
	ID        string
	UserID    string
	Name      string
	AvatarURL string

	CreatedAt time.Time
	UpdatedAt time.Time
}
