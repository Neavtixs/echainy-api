package dto

import (
	"context"

	"github.com/Neavtixs/go-backend-template/internal/apps/domain/entity"
)

type InputRegister struct {
	Ctx      context.Context
	Name     string
	Email    string
	Password string
}

type ResultRegister struct {
	User         entity.User
	Jwt          string
	RefreshToken string
}
