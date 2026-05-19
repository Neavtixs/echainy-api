package dto

type RegisterReq struct {
	Name     string `json:"name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type RegisterRes struct {
	Email string `json:"email"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type LoginRes struct {
	Email string `json:"email"`
}

type MeRes struct {
	ID           string           `json:"id"`
	Email        string           `json:"email"`
	Name         string           `json:"name"`
	AvatarURL    string           `json:"avatar_url"`
	ProviderName string           `json:"provider_name"`
	Workspaces   []MeWorkspaceRes `json:"workspaces"`
}

type MeWorkspaceRes struct {
	ID          string `json:"id"`
	OwnerUserID string `json:"owner_user_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	AvatarURL   string `json:"avatar_url"`
	Role        string `json:"role"`
}

type NewWorkspaceReq struct {
	Name      string `json:"name" validate:"required,max=255"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,max=255"`
}

type NewWorkspaceRes struct {
	ID          string `json:"id"`
	OwnerUserID string `json:"owner_user_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	AvatarURL   string `json:"avatar_url"`
	Role        string `json:"role"`
}

type ListWorkspaceRes struct {
	Workspaces []ListWorkspaceItemRes `json:"workspaces"`
}

type ListWorkspaceItemRes struct {
	ID          string `json:"id"`
	OwnerUserID string `json:"owner_user_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	AvatarURL   string `json:"avatar_url"`
	Role        string `json:"role"`
}
