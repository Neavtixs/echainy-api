package auth

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Neavtixs/echainy-api/internal/apps/domain/entity"
	"github.com/Neavtixs/echainy-api/internal/apps/domain/repository"
	"github.com/Neavtixs/echainy-api/internal/dto"
	"github.com/Neavtixs/echainy-api/internal/errs"
	"github.com/Neavtixs/echainy-api/internal/helper"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	DB                  *sql.DB
	Redis               *redis.Client
	UserRepo            *repository.UserRepo
	UserProfileRepo     *repository.UserProfileRepo
	AuthProviderRepo    *repository.AuthProviderRepo
	WorkspaceRepo       *repository.WorkspaceRepo
	WorkspaceMemberRepo *repository.WorkspaceMemberRepo
}

func NewService(db *sql.DB, redis *redis.Client, userRepo *repository.UserRepo, userProfileRepo *repository.UserProfileRepo, authProviderRepo *repository.AuthProviderRepo, workspaceRepo *repository.WorkspaceRepo, workspaceMemberRepo *repository.WorkspaceMemberRepo) *Service {
	return &Service{
		DB:                  db,
		Redis:               redis,
		UserRepo:            userRepo,
		UserProfileRepo:     userProfileRepo,
		AuthProviderRepo:    authProviderRepo,
		WorkspaceRepo:       workspaceRepo,
		WorkspaceMemberRepo: workspaceMemberRepo,
	}
}

func (s *Service) Register(input *dto.InputRegister) (*dto.ResultRegister, error) {
	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	pass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &entity.User{
		ID:       uuid.NewString(),
		Email:    input.Email,
		Password: string(pass),
	}

	if err := s.UserRepo.Create(tx, input.Ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	userProfile := &entity.UserProfile{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		Name:      input.Name,
		AvatarURL: "",
	}

	if err := s.UserProfileRepo.Create(tx, input.Ctx, userProfile); err != nil {
		return nil, fmt.Errorf("create user profile: %w", err)
	}

	authProvider := &entity.AuthProvider{
		ID:           uuid.NewString(),
		UserID:       user.ID,
		ProviderName: "local",
	}

	if err := s.AuthProviderRepo.Create(tx, input.Ctx, authProvider); err != nil {
		return nil, fmt.Errorf("create auth provider: %w", err)
	}

	workspaceName := fmt.Sprintf("%s workspace starter", input.Name)
	workspace := &entity.Workspace{
		ID:          uuid.NewString(),
		OwnerUserID: user.ID,
		Name:        workspaceName,
		Slug:        helper.GenerateSlug(workspaceName),
		AvatarURL:   "",
	}

	if err := s.WorkspaceRepo.Create(tx, input.Ctx, workspace); err != nil {
		return nil, fmt.Errorf("create workspace: %w", err)
	}

	workspaceMember := &entity.WorkspaceMember{
		ID:          uuid.NewString(),
		WorkspaceID: workspace.ID,
		UserID:      user.ID,
		Role:        "OWNER",
	}

	if err := s.WorkspaceMemberRepo.Create(tx, input.Ctx, workspaceMember); err != nil {
		return nil, fmt.Errorf("create workspace member: %w", err)
	}

	jwt, err := helper.GenerateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken := uuid.NewString()
	key := "refresh_token:" + refreshToken
	if err := s.Redis.Set(
		input.Ctx,
		key,
		user.ID,
		7*24*time.Hour,
	).Err(); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &dto.ResultRegister{
		User:         *user,
		AccessToken:  jwt,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Login(input *dto.InputLogin) (*dto.ResultLogin, error) {
	user := &entity.User{}
	err := s.UserRepo.FindByEmail(s.DB, input.Ctx, input.Email, user)
	if err != nil {
		if err == errs.ErrDataNotFound {
			return nil, errs.ErrInvalidEmailPassword
		}

		return nil, fmt.Errorf("find user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errs.ErrInvalidEmailPassword
	}

	jwt, err := helper.GenerateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken := uuid.NewString()
	key := "refresh_token:" + refreshToken
	if err := s.Redis.Set(
		input.Ctx,
		key,
		user.ID,
		7*24*time.Hour,
	).Err(); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &dto.ResultLogin{
		User:         *user,
		AccessToken:  jwt,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Me(input *dto.InputMe) (*dto.ResultMe, error) {
	userID, ok := input.Ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, errs.ErrInvalidAccessToken
	}

	tx, err := s.DB.BeginTx(input.Ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	user := &entity.User{}
	if err := s.UserRepo.FindByID(s.DB, input.Ctx, userID, user); err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	userProfile := &entity.UserProfile{}
	if err := s.UserProfileRepo.FindByUserID(s.DB, input.Ctx, userID, userProfile); err != nil {
		return nil, fmt.Errorf("find user profile by user id: %w", err)
	}

	authProvider := &entity.AuthProvider{}
	if err := s.AuthProviderRepo.FindByUserID(s.DB, input.Ctx, userID, authProvider); err != nil {
		return nil, fmt.Errorf("find auth provider by user id: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &dto.ResultMe{
		ID:           user.ID,
		Email:        user.Email,
		Name:         userProfile.Name,
		AvatarURL:    userProfile.AvatarURL,
		ProviderName: authProvider.ProviderName,
	}, nil
}

func (s *Service) RefreshAccessToken(input *dto.InputRefreshAccessToken) (*dto.ResultRefreshAccessToken, error) {
	if input.RefreshToken == "" {
		return nil, errs.ErrInvalidRefreshToken
	}

	key := "refresh_token:" + input.RefreshToken
	userID, err := s.Redis.GetDel(input.Ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errs.ErrInvalidRefreshToken
		}

		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	user := &entity.User{}
	if err := s.UserRepo.FindByID(s.DB, input.Ctx, userID, user); err != nil {
		if err == errs.ErrDataNotFound {
			return nil, errs.ErrInvalidRefreshToken
		}

		return nil, fmt.Errorf("find user by id: %w", err)
	}

	accessToken, err := helper.GenerateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken := uuid.NewString()
	newKey := "refresh_token:" + refreshToken
	if err := s.Redis.Set(
		input.Ctx,
		newKey,
		user.ID,
		7*24*time.Hour,
	).Err(); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &dto.ResultRefreshAccessToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Logout(input *dto.InputLogout) error {
	if input.RefreshToken == "" {
		return nil
	}

	key := "refresh_token:" + input.RefreshToken
	if err := s.Redis.Del(input.Ctx, key).Err(); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}

	return nil
}
