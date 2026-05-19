package workspace

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

func (h *Handler) NewHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log = log.WithField("handler", "new_workspace")
	log.WithField("layer", "handler").Info("new workspace request received")

	req := dto.NewWorkspaceReq{}
	c.ShouldBindJSON(&req)

	if err := h.Validate.Struct(req); err != nil {
		log.WithError(err).WithField("layer", "handler").Warn("validation failed")
		c.JSON(http.StatusBadRequest, dto.ResponseWeb[map[string]string]{
			Message: "validation failed",
			Data:    helper.ValidationMsg(err),
		})
		return
	}

	data, err := h.Service.New(&dto.InputNewWorkspace{
		Ctx:       c,
		Name:      req.Name,
		AvatarURL: req.AvatarURL,
	})
	if err != nil {
		log.WithError(err).WithField("layer", "handler").Error("new workspace service failed")

		if errors.Is(err, errs.ErrInvalidAccessToken) {
			c.JSON(http.StatusUnauthorized, dto.ResponseWeb[any]{
				Message: errs.ErrInvalidAccessToken.Error(),
			})
			return
		}

		if errors.Is(err, errs.ErrSlugUsed) {
			c.JSON(http.StatusBadRequest, dto.ResponseWeb[any]{
				Message: "validation failed",
				Data: gin.H{
					"name": errs.ErrSlugUsed.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ResponseWeb[any]{
			Message: errs.ErrInternal.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseWeb[dto.NewWorkspaceRes]{
		Message: "create workspace success",
		Data: dto.NewWorkspaceRes{
			ID:          data.ID,
			OwnerUserID: data.OwnerUserID,
			Name:        data.Name,
			Slug:        data.Slug,
			AvatarURL:   data.AvatarURL,
			Role:        data.Role,
		},
	})
	log.WithField("layer", "handler").Info("new workspace response sent")
}

func (h *Handler) ListHandler(c *gin.Context) {
	log := helper.NewLog(h.Log, c)
	log = log.WithField("handler", "list_workspace")
	log.WithField("layer", "handler").Info("list workspace request received")

	data, err := h.Service.List(&dto.InputListWorkspace{
		Ctx: c,
	})
	if err != nil {
		log.WithError(err).WithField("layer", "handler").Error("list workspace service failed")

		if errors.Is(err, errs.ErrInvalidAccessToken) {
			c.JSON(http.StatusUnauthorized, dto.ResponseWeb[any]{
				Message: errs.ErrInvalidAccessToken.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ResponseWeb[any]{
			Message: errs.ErrInternal.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.ResponseWeb[dto.ListWorkspaceRes]{
		Message: "get workspace list success",
		Data: dto.ListWorkspaceRes{
			Workspaces: mapListWorkspaces(data.Workspaces),
		},
	})
	log.WithField("layer", "handler").Info("list workspace response sent")
}

func mapListWorkspaces(workspaces []dto.ResultListWorkspaceItem) []dto.ListWorkspaceItemRes {
	result := make([]dto.ListWorkspaceItemRes, 0, len(workspaces))
	for _, workspace := range workspaces {
		result = append(result, dto.ListWorkspaceItemRes{
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
