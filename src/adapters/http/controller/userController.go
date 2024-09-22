package controller

import (
	"errors"
	"strings"

	"github.com/Dialosoft/src/adapters/http/request"
	"github.com/Dialosoft/src/adapters/http/response"
	"github.com/Dialosoft/src/adapters/mapper"
	"github.com/Dialosoft/src/domain/services"
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
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "OK", users)
}

func (uc *UserController) GetUserByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	user, err := uc.UserService.GetUserByID(uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
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
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "OK", user)
}

func (uc *UserController) CreateNewUser(c fiber.Ctx) error {
	var req request.UserRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	userDto := mapper.UserRequestToUserDto(&req)

	id, err := uc.UserService.CreateNewUser(*userDto)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return response.ErrConflict(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.StandardCreated(c, "CREATED", id)
}

func (uc *UserController) UpdateUser(c fiber.Ctx) error {
	var req request.UpdateUserRequest

	id := c.Params("id")
	if id == "" {
		return response.ErrEmptyParametersOrArguments(c)
	}

	userUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	if err := c.Bind().Body(&req); err != nil {
		return response.ErrBadRequest(c)
	}

	userDto := mapper.UserUpdateRequestToUserDto(&req)

	err = uc.UserService.UpdateUser(userUUID, *userDto)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "UPDATED", nil)
}

func (uc *UserController) DeleteUser(c fiber.Ctx) error {
	id := c.Params("id")

	userUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	if err = uc.UserService.DeleteUser(userUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "DELETED", nil)
}

func (uc *UserController) RestoreUser(c fiber.Ctx) error {
	id := c.Params("id")

	userUUID, err := uuid.Parse(id)
	if err != nil {
		return response.ErrUUIDParse(c)
	}

	if err = uc.UserService.RestoreUser(userUUID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.ErrNotFound(c)
		}
		return response.ErrInternalServer(c)
	}

	return response.Standard(c, "RESTORED", nil)
}
