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
}

func NewForumController(forumService services.ForumService) *ForumController {
	return &ForumController{ForumService: forumService}
}

func (fc *ForumController) GetAllForums(c fiber.Ctx) error {
	forumsDto, err := fc.ForumService.GetAllForums()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("No forums found", map[string]interface{}{
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving all forums", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Forums retrieved successfully", map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
		"count":  len(forumsDto),
	})

	return response.Standard(c, "OK", forumsDto)
}

func (fc *ForumController) GetForumByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	forumUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	forumDto, err := fc.ForumService.GetForumByID(forumUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Forum not found", map[string]interface{}{
				"forumID": id,
				"route":   c.Path(),
				"method":  c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving forum by ID", map[string]interface{}{
			"forumID": id,
			"route":   c.Path(),
			"method":  c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Forum retrieved successfully", map[string]interface{}{
		"forumID": id,
		"route":   c.Path(),
		"method":  c.Method(),
	})

	return response.Standard(c, "OK", forumDto)
}

func (fc *ForumController) GetForumByName(c fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	forumDto, err := fc.ForumService.GetForumByName(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Forum not found", map[string]interface{}{
				"forumName": name,
				"route":     c.Path(),
				"method":    c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving forum by name", map[string]interface{}{
			"forumName": name,
			"route":     c.Path(),
			"method":    c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Forum retrieved successfully", map[string]interface{}{
		"forumName": name,
		"route":     c.Path(),
		"method":    c.Method(),
	})

	return response.Standard(c, "OK", forumDto)
}

func (fc *ForumController) GetForumsByCategoryIDAndAllowed(c fiber.Ctx) error {
	categoryID := c.Params("categoryID")
	if categoryID == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	categoryUUID, err := uuid.Parse(categoryID)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": categoryID,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
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
		return response.ErrInternalServer(c)
	}

	if forums == nil {
		return response.ErrNotFound(c)
	}

	return response.Standard(c, "OK", forums)
}

func (fc *ForumController) CreateForum(c fiber.Ctx) error {
	var req request.NewForum
	if err := c.Bind().Body(&req); err != nil {
		logger.CaptureError(err, "Failed to parse CreateForum request", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrBadRequest(c)
	}

	err := devconfig.SetDefaultValues(&req)
	if err != nil {
		return response.ErrInternalServer(c)
	}

	forumUUID, err := fc.ForumService.CreateForum(req)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			logger.Warn("Forum creation conflict", map[string]interface{}{
				"request": req,
				"route":   c.Path(),
				"method":  c.Method(),
			})
			return response.ErrConflict(c)
		}
		logger.CaptureError(err, "Error creating new forum", map[string]interface{}{
			"request": req,
			"route":   c.Path(),
			"method":  c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Forum created successfully", map[string]interface{}{
		"forumID": forumUUID.String(),
		"route":   c.Path(),
		"method":  c.Method(),
	})

	return response.StandardCreated(c, "CREATED", fiber.Map{
		"id": forumUUID.String(),
	})
}

func (fc *ForumController) UpdateForum(c fiber.Ctx) error {
	var req request.NewForum

	id := c.Params("id")
	if id == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	forumUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	err = fc.ForumService.UpdateForum(forumUUID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Forum not found for update", map[string]interface{}{
				"forumID": id,
				"route":   c.Path(),
				"method":  c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error updating forum", map[string]interface{}{
			"forumID": id,
			"route":   c.Path(),
			"method":  c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Forum updated successfully", map[string]interface{}{
		"forumID": id,
		"route":   c.Path(),
		"method":  c.Method(),
	})

	return response.Standard(c, "UPDATED", nil)
}

func (fc *ForumController) DeleteForum(c fiber.Ctx) error {
	id := c.Params("id")

	forumUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	if err = fc.ForumService.DeleteForum(forumUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Forum not found for deletion", map[string]interface{}{
				"forumID": id,
				"route":   c.Path(),
				"method":  c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error deleting forum", map[string]interface{}{
			"forumID": id,
			"route":   c.Path(),
			"method":  c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Forum deleted successfully", map[string]interface{}{
		"forumID": id,
		"route":   c.Path(),
		"method":  c.Method(),
	})

	return response.Standard(c, "DELETED", nil)
}

func (fc *ForumController) RestoreForum(c fiber.Ctx) error {
	id := c.Params("id")
	forumUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	if err = fc.ForumService.RestoreForum(forumUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("Forum not found for restoration", map[string]interface{}{
				"forumID": id,
				"route":   c.Path(),
				"method":  c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error restoring forum", map[string]interface{}{
			"forumID": id,
			"route":   c.Path(),
			"method":  c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Forum restored successfully", map[string]interface{}{
		"forumID": id,
		"route":   c.Path(),
		"method":  c.Method(),
	})

	return response.Standard(c, "RESTORED", nil)
}
