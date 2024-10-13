package controller

import (
	"errors"
	"strings"

	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/utils/devconfig"
	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ForumController struct {
	ForumService services.ForumService
	Layer        string
}

func NewForumController(forumService services.ForumService, Layer string) *ForumController {
	return &ForumController{ForumService: forumService, Layer: Layer}
}

func (fc *ForumController) GetAllForums(c fiber.Ctx) error {
	forumsDto, err := fc.ForumService.GetAllForums()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, fc.Layer)
		}
		return response.ErrInternalServer(c, err, forumsDto, fc.Layer)
	}

	return response.Standard(c, "OK", forumsDto)
}

func (fc *ForumController) GetForumByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	forumUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	forumDto, err := fc.ForumService.GetForumByID(forumUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, fc.Layer)
		}
		return response.ErrInternalServer(c, err, forumDto, fc.Layer)
	}

	return response.Standard(c, "OK", forumDto)
}

func (fc *ForumController) GetForumByName(c fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	forumDto, err := fc.ForumService.GetForumByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, fc.Layer)
		}
		return response.ErrInternalServer(c, err, forumDto, fc.Layer)
	}

	return response.Standard(c, "OK", forumDto)
}

func (fc *ForumController) GetForumsByCategoryIDAndAllowed(c fiber.Ctx) error {
	categoryID := c.Params("categoryID")
	if categoryID == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	categoryUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return response.ErrUUIDParse(c, categoryID)
	}

	roleID := c.Locals("roleID")
	roleIDString, ok := roleID.(string)
	if !ok {
		logger.Error("Invalid roleID format in token", map[string]interface{}{
			"roleID": roleID,
			"route":  c.Path(),
		})
		return response.PersonalizedErr(c, "Error in token: claims", fiber.StatusForbidden)
	}

	forums, err := fc.ForumService.GetForumsByCategoryIDAndAllowed(categoryUUID, roleIDString)
	if err != nil {
		return response.ErrInternalServer(c, err, forums, fc.Layer)
	}

	if forums == nil {
		return response.ErrNotFound(c, fc.Layer)
	}

	return response.Standard(c, "OK", forums)
}

func (fc *ForumController) CreateForum(c fiber.Ctx) error {
	var req request.NewForum
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, fc.Layer)
	}

	err := devconfig.SetDefaultValues(&req)
	if err != nil {
		return response.ErrInternalServer(c, err, req, fc.Layer)
	}

	forumUUID, err := fc.ForumService.CreateForum(req)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return response.ErrConflict(c, err, req, fc.Layer)
		}
		return response.ErrInternalServer(c, err, req, fc.Layer)
	}

	return response.StandardCreated(c, "CREATED", fiber.Map{
		"id": forumUUID.String(),
	})
}

func (fc *ForumController) UpdateForum(c fiber.Ctx) error {
	var req request.NewForum

	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	forumUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	err = fc.ForumService.UpdateForum(forumUUID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, fc.Layer)
		}
		return response.ErrInternalServer(c, err, req, fc.Layer)
	}

	return response.Standard(c, "UPDATED", nil)
}

func (fc *ForumController) DeleteForum(c fiber.Ctx) error {
	id := c.Params("id")

	forumUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	if err = fc.ForumService.DeleteForum(forumUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, fc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, fc.Layer)
	}

	return response.Standard(c, "DELETED", nil)
}

func (fc *ForumController) RestoreForum(c fiber.Ctx) error {
	id := c.Params("id")
	forumUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	if err = fc.ForumService.RestoreForum(forumUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, fc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, fc.Layer)
	}

	return response.Standard(c, "RESTORED", nil)
}
