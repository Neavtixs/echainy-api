package auth

import (
	"database/sql"

	"github.com/Neavtixs/go-backend-template/internal/apps/domain/repository"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	DB       *sql.DB
	Redis    *redis.Client
	UserRepo *repository.UserRepo
}

func NewService(db *sql.DB, userRepo *repository.UserRepo, redis *redis.Client) *Service {
	return &Service{
		DB:       db,
		Redis:    redis,
		UserRepo: userRepo,
	}
}

// func (s *Service) Register(input *dto.InputRegister) (*dto.ResultRegister, error) {
// 	tx, err := s.DB.BeginTx(input.Ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer tx.Rollback()

// 	pass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, err
// 	}

// 	user := &entity.User{
// 		ID:       uuid.NewString(),
// 		Email:    input.Email,
// 		Password: string(pass),
// 	}

// 	if err := s.UserRepo.Create(tx, input.Ctx, user); err != nil {
// 		return nil, err
// 	}

// 	userProfile := &entity.UserProfile{
// 		ID:              uuid.NewString(),
// 		UserID:          user.ID,
// 		Name:            input.Name,
// 		Avatar:          "",
// 		Role:            "USER",
// 		ExperienceLimit: 0,
// 		IsGoogle:        false,
// 	}

// 	if err := s.UserProfileRepo.Create(tx, input.Ctx, userProfile); err != nil {
// 		return nil, err
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return nil, err
// 	}

// 	jwt, err := helper.GenerateJWT(user.ID)
// 	if err != nil {
// 		log.Println("redis: ", err)
// 		return nil, err
// 	}

// 	refreshToken := uuid.NewString()
// 	key := "refresh_token:" + refreshToken
// 	if err := s.Redis.Set(
// 		input.Ctx,
// 		key,
// 		user.ID,
// 		7*24*time.Hour,
// 	).Err(); err != nil {
// 		return nil, err
// 	}

// 	return &dto.ResultRegister{
// 		User:         *user,
// 		UserProfile:  *userProfile,
// 		Jwt:          jwt,
// 		RefreshToken: refreshToken,
// 	}, nil
// }
