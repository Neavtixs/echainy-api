package entity

import "time"

type AuthProvider struct {
	ID           string
	UserID       string
	ProviderName string

	CreatedAt time.Time
	UpdatedAt time.Time
}
