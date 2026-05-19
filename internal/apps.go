package internal

import (
	"database/sql"

	"github.com/Neavtixs/echainy-api/internal/apps/domain/repository"
	"github.com/Neavtixs/echainy-api/internal/apps/feature/auth"
	"github.com/Neavtixs/echainy-api/internal/apps/feature/workspace"
	"github.com/Neavtixs/echainy-api/internal/route"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type AppsConfig struct {
	DB       *sql.DB
	App      *gin.Engine
	Redis    *redis.Client
	Validate *validator.Validate
	Log      *logrus.Logger
}

func Apps(a *AppsConfig) {
	userRepo := repository.NewUserRepo()
	userProfileRepo := repository.NewUserProfileRepo()
	authProviderRepo := repository.NewAuthProviderRepo()
	workspaceRepo := repository.NewWorkspaceRepo()
	workspaceMemberRepo := repository.NewWorkspaceMemberRepo()

	authService := auth.NewService(a.DB, a.Redis, userRepo, userProfileRepo, authProviderRepo, workspaceRepo, workspaceMemberRepo)
	authHandler := auth.NewHandler(authService, a.Validate, a.Log)

	workspaceService := workspace.NewService(a.DB, workspaceRepo, workspaceMemberRepo)
	workspaceHandler := workspace.NewHandler(workspaceService, a.Validate, a.Log)

	route.Route{
		App:              a.App,
		AuthHandler:      authHandler,
		WorkspaceHandler: workspaceHandler,
		Log:              a.Log,
	}.SetupRoutes()
}
