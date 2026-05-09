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
