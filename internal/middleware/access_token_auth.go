package middleware

import (
	"net/http"

	"github.com/Neavtixs/go-backend-template/internal/dto"
	"github.com/Neavtixs/go-backend-template/internal/errs"
	"github.com/Neavtixs/go-backend-template/internal/helper"
	"github.com/gin-gonic/gin"
)

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.JSON(401, dto.ResponseWeb[any]{
				Message: errs.ErrInvalidAccessToken.Error(),
			})
			ctx.Abort()
			return
		}

		claims, err := helper.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dto.ResponseWeb[any]{
				Message: errs.ErrInvalidAccessToken.Error(),
			})
			ctx.Abort()
			return
		}

		idVal, ok := claims["user_id"].(string)
		if !ok || idVal == "" {
			ctx.JSON(http.StatusUnauthorized, dto.ResponseWeb[any]{
				Message: errs.ErrInvalidAccessToken.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", idVal)
		ctx.Next()
	}
}
