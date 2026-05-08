package auth

import (
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

// Helper untuk menentukan domain cookie berdasarkan Origin
// func getCookieDomain(c *gin.Context) string {
// 	origin := c.GetHeader("Origin")
// 	log.Println("origin :", origin)
// 	result, _ := url.Parse(origin)

// 	return strings.Split(result.Host, ":")[0]
// }

// func (h *Handler) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
// 	c.SetCookieData(&http.Cookie{
// 		Name:     "refresh_token",
// 		Value:    refreshToken,
// 		HttpOnly: true,
// 		SameSite: http.SameSiteNoneMode,
// 		Secure:   true,
// 		Path:     "/api/auth/refresh",
// 	})
// 	c.SetCookieData(&http.Cookie{
// 		Name:     "access_token",
// 		Value:    accessToken,
// 		HttpOnly: true,
// 		SameSite: http.SameSiteNoneMode,
// 		Secure:   true,
// 	})
// }

// func (h *Handler) clearAuthCookies(c *gin.Context) {
// 	c.SetCookieData(&http.Cookie{
// 		Name:     "refresh_token",
// 		Value:    "",
// 		HttpOnly: true,
// 		SameSite: http.SameSiteNoneMode,
// 		Path:     "/api/auth/refresh",
// 		MaxAge:   -1,
// 		Secure:   true,
// 	})
// 	c.SetCookieData(&http.Cookie{
// 		Name:     "access_token",
// 		Value:    "",
// 		HttpOnly: true,
// 		SameSite: http.SameSiteNoneMode,
// 		MaxAge:   -1,
// 		Secure:   true,
// 	})
// }

// func (h *Handler) RegisterHandler(c *gin.Context) {
// 	log := helper.NewLog(h.Log, c)
// 	log.WithField("layer", "handler").Info("register user start")

// 	req := dto.RegisterReq{}
// 	c.ShouldBindJSON(&req)

// 	if err := h.Validate.Struct(req); err != nil {
// 		log.WithField("layer", "handler").Error(err)
// 		c.JSON(400, dto.ResponseWeb[map[string]string]{
// 			Message: "validation failed",
// 			Data:    helper.ValidationMsg(err),
// 		})
// 		return

// 	}

// 	data, err := h.Service.Register(&dto.InputRegister{
// 		Ctx:      c,
// 		Name:     req.Name,
// 		Email:    req.Email,
// 		Password: req.Password,
// 	})
// 	if err != nil {
// 		log.WithField("layer", "register service").Error(err)

// 		if errors.Is(err, errs.ErrEmailUsed) {
// 			c.JSON(400, dto.ResponseWeb[any]{
// 				Message: "validation vailed",
// 				Data: gin.H{
// 					"email": err.Error(),
// 				},
// 			})
// 			return
// 		}

// 		c.JSON(http.StatusInternalServerError, dto.ResponseWeb[any]{
// 			Message: errs.ErrInternal.Error(),
// 		})
// 		return
// 	}

// 	h.setAuthCookies(c, data.Jwt, data.RefreshToken)

// 	c.JSON(http.StatusOK, dto.ResponseWeb[dto.RegisterRes]{
// 		Message: "register user success",
// 		Data: dto.RegisterRes{
// 			Email: data.User.Email,
// 		},
// 	})

// }
