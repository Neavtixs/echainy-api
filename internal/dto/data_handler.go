package dto

type RegisterReq struct {
	Name     string `json:"name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type RegisterRes struct {
	Email string `json:"email"`
}
