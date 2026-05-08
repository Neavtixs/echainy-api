package errs

import "errors"

var (
	ErrEmailUsed              = errors.New("email already used")
	ErrUserIDUsed             = errors.New("user_id already used")
	ErrUserIDNotFound         = errors.New("user_id is not found")
	ErrFailedCreateData       = errors.New("failed create database")
	ErrInternal               = errors.New("internal server error")
	ErrInvalidType            = errors.New("invalid type")
	ErrDataNotFound           = errors.New("data is not found")
	ErrInvalidEmailPassword   = errors.New("email or password is invalid")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
	ErrInvalidAccessToken     = errors.New("invalid access token")
	ErrInvalidStateGoogle     = errors.New("invalid state google")
	ErrEmailIsNotGoogleLogin  = errors.New("this email already used without google")
	ErrEmailIsGoogleLogin     = errors.New("this email was used with google login")
	ErrInvalidRole            = errors.New("role must be USER or ADMIN")
	ErrInvalidExperienceLimit = errors.New("experience_limit must be greater than or equal to 0")
	ErrCannotChangeOwnRole    = errors.New("admin cannot change their own role")
	ErrInvalidImageFile       = errors.New("invalid image file, allowed: jpg, jpeg, png, webp")
	ErrImageTooLarge          = errors.New("image size must be less than or equal to 5MB")

	ErrExperienceLimitReached = errors.New("experience limist user is reached")
)
