package dto

import (
	"context"

	"github.com/Neavtixs/echainy-api/internal/apps/domain/entity"
)

type InputRegister struct {
	Ctx      context.Context
	Name     string
	Email    string
	Password string
}

type ResultRegister struct {
	User         entity.User
	AccessToken  string
	RefreshToken string
}

type InputLogin struct {
	Ctx      context.Context
	Email    string
	Password string
}

type ResultLogin struct {
	User         entity.User
	AccessToken  string
	RefreshToken string
}

type InputMe struct {
	Ctx context.Context
}

type ResultMe struct {
	ID           string
	Email        string
	Name         string
	AvatarURL    string
	ProviderName string
	Workspaces   []ResultMeWorkspace
}

type ResultMeWorkspace struct {
	ID          string
	OwnerUserID string
	Name        string
	Slug        string
	AvatarURL   string
	Role        string
}

type InputRefreshAccessToken struct {
	Ctx          context.Context
	RefreshToken string
}

type ResultRefreshAccessToken struct {
	AccessToken  string
	RefreshToken string
}

type InputLogout struct {
	Ctx          context.Context
	RefreshToken string
}

type InputNewWorkspace struct {
	Ctx       context.Context
	Name      string
	AvatarURL string
}

type ResultNewWorkspace struct {
	ID          string
	OwnerUserID string
	Name        string
	Slug        string
	AvatarURL   string
	Role        string
}

type InputListWorkspace struct {
	Ctx context.Context
}

type ResultListWorkspace struct {
	Workspaces []ResultListWorkspaceItem
}

type ResultListWorkspaceItem struct {
	ID          string
	OwnerUserID string
	Name        string
	Slug        string
	AvatarURL   string
	Role        string
}
