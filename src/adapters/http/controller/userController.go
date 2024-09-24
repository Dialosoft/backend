package controller

import (
	"errors"
	"strings"

	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/domain/services"
	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (uc *UserController) GetAllUsers(c fiber.Ctx) error {
	users, err := uc.UserService.GetAllUsers()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("No users found", map[string]interface{}{
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving all users", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("Users retrieved successfully", map[string]interface{}{
		"count":  len(users),
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "OK", users)
}

func (uc *UserController) GetUserByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	userUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	user, err := uc.UserService.GetUserByID(userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("User not found", map[string]interface{}{
				"userID": id,
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving user by ID", map[string]interface{}{
			"userID": id,
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("User retrieved successfully", map[string]interface{}{
		"userID": id,
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "OK", user)
}

func (uc *UserController) GetUserByUsername(c fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	user, err := uc.UserService.GetUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("User not found by username", map[string]interface{}{
				"username": username,
				"route":    c.Path(),
				"method":   c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error retrieving user by username", map[string]interface{}{
			"username": username,
			"route":    c.Path(),
			"method":   c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("User retrieved successfully by username", map[string]interface{}{
		"username": username,
		"route":    c.Path(),
		"method":   c.Method(),
	})

	return response.Standard(c, "OK", user)
}

func (uc *UserController) CreateNewUser(c fiber.Ctx) error {
	var req request.UserRequest
	if err := c.Bind().Body(&req); err != nil {
		logger.CaptureError(err, "Failed to bind request for creating new user", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrBadRequest(c)
	}

	userDto := mapper.UserRequestToUserDto(&req)

	id, err := uc.UserService.CreateNewUser(*userDto)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			logger.Warn("Duplicate user detected", map[string]interface{}{
				"userDto": userDto,
				"route":   c.Path(),
				"method":  c.Method(),
			})
			return response.ErrConflict(c)
		}
		logger.CaptureError(err, "Error creating new user", map[string]interface{}{
			"userDto": userDto,
			"route":   c.Path(),
			"method":  c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("User created successfully", map[string]interface{}{
		"userID": id.String(),
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.StandardCreated(c, "CREATED", id)
}

func (uc *UserController) UpdateUser(c fiber.Ctx) error {
	var req request.NewUser

	id := c.Params("id")
	if id == "" {
		logger.Error("Empty parameters or arguments", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrEmptyParametersOrArguments(c)
	}

	userUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	if err := c.Bind().Body(&req); err != nil {
		logger.CaptureError(err, "Failed to bind request for updating user", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrBadRequest(c)
	}

	err = uc.UserService.UpdateUser(userUUID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("User not found for update", map[string]interface{}{
				"userID": id,
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error updating user", map[string]interface{}{
			"userID": id,
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("User updated successfully", map[string]interface{}{
		"userID": id,
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "UPDATED", nil)
}

func (uc *UserController) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")

	userUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	if err = uc.UserService.DeleteUser(userUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("User not found for deletion", map[string]interface{}{
				"userID": id,
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error deleting user", map[string]interface{}{
			"userID": id,
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("User deleted successfully", map[string]interface{}{
		"userID": id,
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "DELETED", nil)
}

func (uc *UserController) RestoreUser(c fiber.Ctx) error {
	id := c.Params("id")

	userUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error("Invalid UUID format", map[string]interface{}{
			"provided-id": id,
			"route":       c.Path(),
			"method":      c.Method(),
		})
		return response.ErrUUIDParse(c)
	}

	if err = uc.UserService.RestoreUser(userUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Warn("User not found for restoration", map[string]interface{}{
				"userID": id,
				"route":  c.Path(),
				"method": c.Method(),
			})
			return response.ErrNotFound(c)
		}
		logger.CaptureError(err, "Error restoring user", map[string]interface{}{
			"userID": id,
			"route":  c.Path(),
			"method": c.Method(),
		})
		return response.ErrInternalServer(c)
	}

	logger.Info("User restored successfully", map[string]interface{}{
		"userID": id,
		"route":  c.Path(),
		"method": c.Method(),
	})

	return response.Standard(c, "RESTORED", nil)
}
