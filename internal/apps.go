package internal

import (
	"database/sql"

	"github.com/Neavtixs/echainy-api/internal/route"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type AppsConfig struct {
	DB		*sql.DB
	App		*gin.Engine
	Redis		*redis.Client
	Validate	*validator.Validate
	Log		*logrus.Logger
}

func Apps(a *AppsConfig) {
	// userRepo := repository.NewUserRepo()

	// authService := auth.NewService(a.DB, userRepo, userProfileRepo, a.Redis)
	// authHandler := auth.NewHandler(authService, a.Validate, a.Log)
	route.Route{
		App: a.App,
		// AuthHandler: authHandler,
	}.SetupRoutes()
}
