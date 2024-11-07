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
	Layer       string
}

func NewUserController(userService services.UserService, Layer string) *UserController {
	return &UserController{UserService: userService, Layer: Layer}
}

func (uc *UserController) GetAllUsers(c fiber.Ctx) error {
	users, err := uc.UserService.GetAllUsers()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, uc.Layer)
		}
		return response.ErrInternalServer(c, err, users, uc.Layer)
	}

	return response.Standard(c, "OK", users)
}

func (uc *UserController) GetUserByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	userUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	user, err := uc.UserService.GetUserByID(userUUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, uc.Layer)
		}
		return response.ErrInternalServer(c, err, user, uc.Layer)
	}

	return response.Standard(c, "OK", user)
}

func (uc *UserController) GetUserByUsername(c fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	user, err := uc.UserService.GetUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, uc.Layer)
		}
		return response.ErrInternalServer(c, err, user, uc.Layer)
	}

	return response.Standard(c, "OK", user)
}

func (uc *UserController) CreateNewUser(c fiber.Ctx) error {
	var req request.UserRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, uc.Layer)
	}

	userDto := mapper.UserRequestToUserDto(&req)

	id, err := uc.UserService.CreateNewUser(*userDto)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return response.ErrConflict(c, err, userDto, uc.Layer)
		}
		return response.ErrInternalServer(c, err, userDto, uc.Layer)
	}

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
		return response.ErrUUIDParse(c, id)
	}

	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, uc.Layer)
	}

	err = uc.UserService.UpdateUser(userUUID, req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, uc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, uc.Layer)
	}

	return response.Standard(c, "UPDATED", nil)
}

func (uc *UserController) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")

	userUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	if err = uc.UserService.DeleteUser(userUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, uc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, uc.Layer)
	}

	return response.Standard(c, "DELETED", nil)
}

func (uc *UserController) RestoreUser(c fiber.Ctx) error {
	id := c.Params("id")

	userUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	if err = uc.UserService.RestoreUser(userUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c, uc.Layer)
		}
		return response.ErrInternalServer(c, err, nil, uc.Layer)
	}

	return response.Standard(c, "RESTORED", nil)
}

func (uc *UserController) ChangeUserAvatar(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	userUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c, id)
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return response.ErrBadRequest(c, string(c.Body()), err, uc.Layer)
	}

	if fileHeader.Size > 5*1024*1024 {
		logger.Warn("File size too large", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
			"size":   fileHeader.Size,
		})
		return response.PersonalizedErr(c, "File size exceeds the 5MB limit", fiber.StatusRequestEntityTooLarge)
	}

	if fileHeader.Header.Get("Content-Type") != "image/png" && fileHeader.Header.Get("Content-Type") != "image/jpeg" {
		logger.Warn("Invalid file type", map[string]interface{}{
			"route":  c.Path(),
			"method": c.Method(),
			"type":   fileHeader.Header.Get("Content-Type"),
		})
		return response.PersonalizedErr(c, "Only PNG and JPEG formats are allowed", fiber.StatusUnsupportedMediaType)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return response.ErrInternalServer(c, err, nil, uc.Layer)
	}
	defer file.Close()

	err = uc.UserService.ProcessAvatar(userUUID, fileHeader, file)
	if err != nil {
		logger.Error("Failed to process avatar upload", map[string]interface{}{
			"user_id": userUUID.String(),
			"route":   c.Path(),
			"method":  c.Method(),
			"error":   err.Error(),
		})
		return response.ErrInternalServer(c, err, nil, uc.Layer)
	}

	return response.Standard(c, "Avatar uploaded successfully", nil)
}
