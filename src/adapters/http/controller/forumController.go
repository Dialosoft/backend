package controller

import (
	"errors"
	"strings"

	"github.com/Dialosoft/src/adapters/dto"
	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ForumController struct {
	ForumService services.ForumService
}

func NewForumController(forumService services.ForumService) *ForumController {
	return &ForumController{ForumService: forumService}
}

func (fc *ForumController) GetAllForums(c fiber.Ctx) error {
	forumsDto, err := fc.ForumService.GetAllForums()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.Standard(c, "OK", forumsDto))
}

func (fc *ForumController) GetForumByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Error("empty parameters or arguments")
		return response.ErrEmptyParametersOrArguments(c)
	}

	forumUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	forumDto, err := fc.ForumService.GetForumByID(forumUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.Standard(c, "OK", forumDto))
}

func (fc *ForumController) GetForumByName(c fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		logger.Error("empty parameters or arguments")
		return response.ErrEmptyParametersOrArguments(c)
	}

	forumDto, err := fc.ForumService.GetForumByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.Standard(c, "OK", forumDto))
}

func (fc *ForumController) CreateForum(c fiber.Ctx) error {
	var req request.CreateForum
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	forumDto := dto.ForumDto{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
	}

	forumUUID, err := fc.ForumService.CreateForum(forumDto)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return response.ErrConflict(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.StandardCreated(c, "CREATED", fiber.Map{
		"id": forumUUID.String(),
	}))
}

func (fc *ForumController) UpdateForum(c fiber.Ctx) error {
	var req request.UpdateForum

	id := c.Params("id")
	if id == "" {
		return response.ErrBadRequest(c)
	}

	forumUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	forumDto := dto.ForumDto{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
	}

	err = fc.ForumService.UpdateForum(forumUUID, forumDto)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.Standard(c, "UPDATED", nil))
}

func (fc *ForumController) DeleteForum(c fiber.Ctx) error {
	id := c.Params("id")

	forumUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	if err = fc.ForumService.DeleteForum(forumUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.Standard(c, "DELETED", nil))
}

func (fc *ForumController) RestoreForum(c fiber.Ctx) error {
	id := c.Params("id")
	forumUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	if err = fc.ForumService.RestoreForum(forumUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return c.Status(fiber.StatusOK).JSON(response.Standard(c, "RESTORED", nil))
}
