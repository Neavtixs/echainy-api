package auth

import (
	"errors"
	"net/http"

	"github.com/Neavtixs/echainy-api/internal/dto"
	"github.com/Neavtixs/echainy-api/internal/errs"
	"github.com/Neavtixs/echainy-api/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Service  *Service
	Validate *validator.Validate
	Log      *logrus.Logger
}

func NewHandler(service *Service, validate *validator.Validate, log *logrus.Logger) *Handler {
	return &Handler{
		Service:  service,
		Validate: validate,
		Log:      log,
	}

}

func (h *Handler) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	c.SetCookieData(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Path:     "/api/auth/refresh",
	})
	c.SetCookieData(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
}

func (h *Handler) clearAuthCookies(c *gin.Context) {
	c.SetCookieData(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/api/auth/refresh",
		MaxAge:   -1,
		Secure:   true,
	})
	c.SetCookieData(&http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   -1,
		Secure:   true,
	})
}

func (h *Handler) RegisterHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log = log.WithField("handler", "register")
	log.WithField("layer", "handler").Info("register request received")

	req := dto.RegisterReq{}
	c.ShouldBindJSON(&req)

	if err := h.Validate.Struct(req); err != nil {
		log.WithError(err).WithField("layer", "handler").Warn("validation failed")
		c.JSON(http.StatusBadRequest, dto.ResponseWeb[map[string]string]{
			Message: "validation failed",
			Data:    helper.ValidationMsg(err),
		})
		return

	}

	data, err := h.Service.Register(&dto.InputRegister{
		Ctx:      c,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.WithError(err).WithField("layer", "handler").Error("register service failed")

		if errors.Is(err, errs.ErrEmailUsed) {
			c.JSON(http.StatusBadRequest, dto.ResponseWeb[any]{
				Message: "validation failed",
				Data: gin.H{
					"email": errs.ErrEmailUsed.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ResponseWeb[any]{
			Message: errs.ErrInternal.Error(),
		})
		return
	}

	h.setAuthCookies(c, data.AccessToken, data.RefreshToken)
	log.WithField("layer", "handler").Info("auth cookies set")

	c.JSON(http.StatusOK, dto.ResponseWeb[dto.RegisterRes]{
		Message: "register user success",
		Data: dto.RegisterRes{
			Email: data.User.Email,
		},
	})
	log.WithField("layer", "handler").Info("register response sent")

}

func (h *Handler) LoginHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log = log.WithField("handler", "login")
	log.WithField("layer", "handler").Info("login request received")

	req := dto.LoginReq{}
	c.ShouldBindJSON(&req)

	if err := h.Validate.Struct(req); err != nil {
		log.WithError(err).WithField("layer", "handler").Warn("validation failed")
		c.JSON(http.StatusBadRequest, dto.ResponseWeb[map[string]string]{
			Message: "validation failed",
			Data:    helper.ValidationMsg(err),
		})
		return
	}

	data, err := h.Service.Login(&dto.InputLogin{
		Ctx:      c,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.WithError(err).WithField("layer", "handler").Error("login service failed")

		if errors.Is(err, errs.ErrInvalidEmailPassword) {
			c.JSON(http.StatusUnauthorized, dto.ResponseWeb[any]{
				Message: errs.ErrInvalidEmailPassword.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ResponseWeb[any]{
			Message: errs.ErrInternal.Error(),
		})
		return
	}

	h.setAuthCookies(c, data.AccessToken, data.RefreshToken)
	log.WithField("layer", "handler").Info("auth cookies set")

	c.JSON(http.StatusOK, dto.ResponseWeb[dto.LoginRes]{
		Message: "login user success",
		Data: dto.LoginRes{
			Email: data.User.Email,
		},
	})
	log.WithField("layer", "handler").Info("login response sent")
}

func (h *Handler) MeHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log = log.WithField("handler", "me")
	log.WithField("layer", "handler").Info("me request received")

	data, err := h.Service.Me(&dto.InputMe{
		Ctx: c,
	})
	if err != nil {
		log.WithError(err).WithField("layer", "handler").Error("me service failed")

		if errors.Is(err, errs.ErrInvalidAccessToken) {
			c.JSON(http.StatusUnauthorized, dto.ResponseWeb[any]{
				Message: errs.ErrInvalidAccessToken.Error(),
			})
			return
		}

		if errors.Is(err, errs.ErrDataNotFound) {
			c.JSON(http.StatusNotFound, dto.ResponseWeb[any]{
				Message: errs.ErrDataNotFound.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ResponseWeb[any]{
			Message: errs.ErrInternal.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseWeb[dto.MeRes]{
		Message: "get me success",
		Data: dto.MeRes{
			ID:           data.ID,
			Email:        data.Email,
			Name:         data.Name,
			AvatarURL:    data.AvatarURL,
			ProviderName: data.ProviderName,
			Workspaces:   mapMeWorkspaces(data.Workspaces),
		},
	})
	log.WithField("layer", "handler").Info("me response sent")
}

func mapMeWorkspaces(workspaces []dto.ResultMeWorkspace) []dto.MeWorkspaceRes {
	result := make([]dto.MeWorkspaceRes, 0, len(workspaces))
	for _, workspace := range workspaces {
		result = append(result, dto.MeWorkspaceRes{
			ID:          workspace.ID,
			OwnerUserID: workspace.OwnerUserID,
			Name:        workspace.Name,
			Slug:        workspace.Slug,
			AvatarURL:   workspace.AvatarURL,
			Role:        workspace.Role,
		})
	}

	return result
}

func (h *Handler) RefreshAccessTokenHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log = log.WithField("handler", "refresh_access_token")
	log.WithField("layer", "handler").Info("refresh access token request received")

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.WithError(err).WithField("layer", "handler").Warn("refresh token cookie not found")
		c.JSON(http.StatusUnauthorized, dto.ResponseWeb[any]{
			Message: errs.ErrInvalidRefreshToken.Error(),
		})
		return
	}

	data, err := h.Service.RefreshAccessToken(&dto.InputRefreshAccessToken{
		Ctx:          c,
		RefreshToken: refreshToken,
	})
	if err != nil {
		log.WithError(err).WithField("layer", "handler").Error("refresh access token service failed")

		if errors.Is(err, errs.ErrInvalidRefreshToken) {
			c.JSON(http.StatusUnauthorized, dto.ResponseWeb[any]{
				Message: errs.ErrInvalidRefreshToken.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ResponseWeb[any]{
			Message: errs.ErrInternal.Error(),
		})
		return
	}

	h.setAuthCookies(c, data.AccessToken, data.RefreshToken)
	log.WithField("layer", "handler").Info("auth cookies refreshed")

	c.JSON(http.StatusOK, dto.ResponseWeb[any]{
		Message: "refresh access token success",
	})
	log.WithField("layer", "handler").Info("refresh access token response sent")
}

func (h *Handler) LogoutHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log = log.WithField("handler", "logout")
	log.WithField("layer", "handler").Info("logout request received")

	refreshToken, _ := c.Cookie("refresh_token")
	if err := h.Service.Logout(&dto.InputLogout{
		Ctx:          c,
		RefreshToken: refreshToken,
	}); err != nil {
		log.WithError(err).WithField("layer", "handler").Error("logout service failed")
		c.JSON(http.StatusInternalServerError, dto.ResponseWeb[any]{
			Message: errs.ErrInternal.Error(),
		})
		return
	}

	h.clearAuthCookies(c)
	log.WithField("layer", "handler").Info("auth cookies cleared")

	c.JSON(http.StatusOK, dto.ResponseWeb[any]{
		Message: "logout success",
	})
	log.WithField("layer", "handler").Info("logout response sent")
}
