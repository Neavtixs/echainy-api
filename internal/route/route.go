package route

import (
	"net/http"
	"os"
	"strings"

	"github.com/Neavtixs/echainy-api/internal/apps/feature/auth"
	"github.com/Neavtixs/echainy-api/internal/dto"
	"github.com/Neavtixs/echainy-api/internal/helper"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Route struct {
	App         *gin.Engine
	AuthHandler *auth.Handler
	Log         *logrus.Logger
}

func (r Route) SetupRoutes() {
	r.App.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(os.Getenv("CORS_ALLOW_ORIGIN"), ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	r.App.Use(helper.RequestIDMiddleware())
	r.App.Use(helper.RequestLogger(r.Log))

	r.App.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.ResponseWeb[any]{
			Message: "route not found",
		})
	})

	api := r.App.Group("/api")

	api.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.ResponseWeb[any]{
			Message: "this your root api",
		})
	})

	api.Static("/uploads", os.Getenv("UPLOAD_DIR"))

	// ========================
	// PUBLIC ROUTES
	// ========================

	public := api.Group("")
	{
		// // auth public
		public.POST("/auth/register", r.AuthHandler.RegisterHandler)
	}
	// public.POST("/auth/login", r.AuthHandler.LoginHandler)
	// public.POST("/auth/refresh", r.AuthHandler.RefreshAccessTokenHandler)

	// public.GET("/auth/google/login", r.AuthHandler.GoogleRedirectHandler)
	// public.GET("/auth/google/callback", r.AuthHandler.GoogleCallbackHandler)
	// public.POST("/auth/logout", r.AuthHandler.LogoutHandler)

	// // product
	// public.GET("/product/health", r.ProductHandler.HealthHandler)
	// public.GET("/product", r.ProductHandler.ProductListHandler)
	// public.GET("/product/category", r.ProductHandler.CategoryListHandler)
	// public.GET("/product/:id", r.ProductHandler.ProductDetailHandler)

	// // donation
	// public.GET("/donation/notification", r.DonationHandler.GetNotificationHandler)

	// // webhook (external service)
	// public.POST("/webhook/:roblox_experience_id/saweria", r.DonationHandler.SaweriaWebhookHandler)
	// public.POST("/webhook/:roblox_experience_id/bagibagi", r.DonationHandler.BagibagiWebhookHandler)
	// public.POST("/webhook/:roblox_experience_id/sociabuzz", r.DonationHandler.SociabuzzWebhookHandler)

	// // leaderboard
	// public.GET("/leaderboard/saweria", r.DonationHandler.LeaderboardSaweriaHandler)
	// public.GET("/leaderboard/bagibagi", r.DonationHandler.LeaderboardBagibagiHandler)
	// public.GET("/leaderboard/sociabuzz", r.DonationHandler.LeaderboardSociabuzzHandler)
	// public.GET("/leaderboard/map", r.DonationHandler.GetMapLeaderboardHandler)

	// ========================
	// AUTHENTICATED USER
	// ========================
	// user := api.Group("", middleware.Authorization())
	// {
	// 	user.GET("/auth/me", r.AuthHandler.MeHandler)

	// 	// experience
	// 	user.POST("/experience", r.ExperienceHandler.CreateExperienceHandler)
	// 	user.GET("/experience/my-list", r.ExperienceHandler.MyListHandler)
	// 	user.GET("/experience/:id", r.ExperienceHandler.ShowDetailHandler)
	// 	user.POST("/experience/activated/:id", r.ExperienceHandler.SetActivationStatusHandler)

	// 	// donation
	// 	user.GET("/donation/history", r.DonationHandler.HistoryHandler)
	// 	user.GET("/donation/donatur-leaderboard", r.DonationHandler.DonaturLeaderboardHandler)
	// }

	// // ========================
	// // ADMIN ROUTES
	// // ========================
	// admin := api.Group("/admin", middleware.Authorization() /* + middleware.AdminOnly() */)
	// {
	// 	// experience
	// 	admin.GET("/experience", r.ExperienceHandler.AdminAllListHandler)
	// 	admin.GET("/experience/:id", r.ExperienceHandler.AdminShowDetailHandler)
	// 	admin.POST("/experience/activated/:id", r.ExperienceHandler.AdminSetActivationStatusHandler)

	// 	// users
	// 	admin.GET("/users", r.ProfileHandler.AdminListUsersHandler)
	// 	admin.GET("/users/:id", r.ProfileHandler.AdminDetailUsersHandler)
	// 	admin.POST("/users/change-role/:id", r.ProfileHandler.AdminChangeRoleUserHandler)
	// 	admin.POST("/users/limit-experience/:id", r.ProfileHandler.AdminChangeLimitExperienceHandler)

	// 	// product
	// 	admin.POST("/product/category", r.ProductHandler.AdminCreateCategoryHandler)
	// 	admin.POST("/product", r.ProductHandler.AdminCreateProductHandler)
	// 	admin.PUT("/product/:id", r.ProductHandler.AdminUpdateProductHandler)
	// 	admin.DELETE("/product/:id", r.ProductHandler.AdminRemoveProductHandler)
	// 	admin.DELETE("/product/category/:id", r.ProductHandler.AdminRemoveCategoryHandler)
	// }
}
