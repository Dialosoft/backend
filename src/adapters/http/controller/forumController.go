package controller

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/utils/devconfig"
	"github.com/Dialosoft/src/pkg/utils/logger"
)

type ForumController struct {
	ForumService services.ForumService
	RoleService  services.RoleService
	Layer        string
}

func NewForumController(forumService services.ForumService, roleService services.RoleService, Layer string) *ForumController {
	return &ForumController{ForumService: forumService, RoleService: roleService, Layer: Layer}
}

// GetAllForums retrieves all forums.
// @Summary Get all forums
// @Description Retrieves a list of all forums.
// @Tags Forums
// @Produce json
// @Success 200 {object} []dto.ForumDto "OK"
// @Failure 404 {object} response.StandardError "NOT FOUND"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /forums/protected/get-all-forums [get]
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

// GetForumByID retrieves a forum by its ID.
// @Summary Get forum by ID
// @Description Retrieves a specific forum by its unique ID.
// @Tags Forums
// @Param id path string true "Forum ID"
// @Produce json
// @Success 200 {object} dto.ForumDto "OK"
// @Failure 400 {object} response.StandardError "BAD REQUEST"
// @Failure 404 {object} response.StandardError "NOT FOUND"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /forums/protected/get-forum-by-id/{id} [get]
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

// GetForumByName retrieves a forum by its name.
// @Summary Get forum by name
// @Description Retrieves a specific forum by its unique name.
// @Tags Forums
// @Param name path string true "Forum Name"
// @Produce json
// @Success 200 {object} dto.ForumDto "OK"
// @Failure 400 {object} response.StandardError "BAD REQUEST"
// @Failure 404 {object} response.StandardError "NOT FOUND"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /forums/protected/get-forum-by-name/{name} [get]
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

// GetForumsByCategoryIDAndAllowed retrieves forums allowed by category ID.
// @Summary Get forums by category ID
// @Description Retrieves forums by category ID if allowed for the role, due to it's a public endpoint, in case the user is not authenticated, it will return all forums.
// @Tags Forums
// @Param categoryID path string true "Category ID"
// @Produce json
// @Success 200 {object} []dto.ForumDto "OK"
// @Failure 400 {object} response.StandardError "BAD REQUEST"
// @Failure 403 {object} response.StandardError "FORBIDDEN"
// @Failure 404 {object} response.StandardError "NOT FOUND"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /forums/get-forums-by-category-id/{categoryID} [get]
func (fc *ForumController) GetForumsByCategoryIDAndAllowed(c fiber.Ctx) error {
	categoryID := c.Params("categoryID")
	if categoryID == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	categoryUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return response.ErrUUIDParse(c, categoryID)
	}

	// Predetermined value for non -authenticated users
	var defaultUserRole string = "anonymous"

	if roleIDFromContext := c.Locals("roleID"); roleIDFromContext != nil {
		roleIDString, ok := roleIDFromContext.(string)
		if ok {

			roleIdFromContextUUID, err := uuid.Parse(roleIDString)
			
			if err != nil {
				return response.ErrUUIDParse(c, roleIDString)
			}

			roleFromContext, err := fc.RoleService.GetRoleByID(roleIdFromContextUUID)

			if err != nil {
				return response.ErrInternalServer(c, err, roleFromContext, fc.Layer)
			}

			logger.Info("roleId exists in context", map[string]interface{}{
				"roleID": roleIDFromContext,
				"route":  c.Path(),
			})
			defaultUserRole = roleFromContext.RoleType
		}
	}

	forums, err := fc.ForumService.GetForumsByCategoryIDAndAllowed(categoryUUID, defaultUserRole)
	if err != nil {
		return response.ErrInternalServer(c, err, forums, fc.Layer)
	}

	if forums == nil {
		return response.ErrNotFound(c, fc.Layer)
	}

	return response.Standard(c, "OK", forums)
}

// CreateForum creates a new forum.
// @Summary Create a new forum
// @Description Creates a new forum with the specified details.
// @Tags Forums
// @Accept json
// @Produce json
// @Param forum body request.NewForum true "New Forum Data"
// @Success 201 {object} string "CREATED"
// @Failure 400 {object} response.StandardError "BAD REQUEST"
// @Failure 409 {object} response.StandardError "CONFLICT"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /forums/protected/create-new-forum [post]
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

// UpdateForum updates an existing forum.
// @Summary Update a forum
// @Description Updates an existing forum with the specified ID and details.
// @Tags Forums
// @Accept json
// @Produce json
// @Param id path string true "Forum ID"
// @Param forum body request.NewForum true "Updated Forum Data"
// @Success 200 {object} response.StandardResponse "UPDATED"
// @Failure 400 {object} response.StandardError "BAD REQUEST"
// @Failure 404 {object} response.StandardError "NOT FOUND"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /forums/protected/update-forum/{id} [put]
func (fc *ForumController) UpdateForum(c fiber.Ctx) error {
	var req request.NewForum

	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, fc.Layer)
	}

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

// DeleteForum deletes a forum by ID.
// @Summary Delete a forum
// @Description Deletes a forum with the specified ID.
// @Tags Forums
// @Param id path string true "Forum ID"
// @Produce json
// @Success 200 {object} response.StandardResponse "DELETED"
// @Failure 400 {object} response.StandardError "BAD REQUEST"
// @Failure 404 {object} response.StandardError "NOT FOUND"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /forums/protected/delete-forum/{id} [delete]
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

// RestoreForum restores a deleted forum by ID.
// @Summary Restore a forum
// @Description Restores a deleted forum with the specified ID.
// @Tags Forums
// @Param id path string true "Forum ID"
// @Produce json
// @Success 200 {object} response.StandardResponse "RESTORED"
// @Failure 400 {object} response.StandardError "BAD REQUEST"
// @Failure 404 {object} response.StandardError "NOT FOUND"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /forums/protected/restore-forum/{id} [put]
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
