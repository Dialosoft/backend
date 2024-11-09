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

// GetAllUsers retrieves all users in the system.
//
// @Summary Get all users
// @Description Retrieves a list of all registered users.
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} response.StandardResponse "List of users"
// @Failure 404 {object} response.StandardError "NOT FOUND - No users found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router /users/get-all-users [get]
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

// GetUserByID retrieves a user by their unique ID.
//
// @Summary Get user by ID
// @Description Retrieves user details based on their unique identifier.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.StandardResponse "User data"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid ID format"
// @Failure 404 {object} response.StandardError "NOT FOUND - User not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router /users/get-user-by-id/{id} [get]
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

// GetUserByUsername retrieves a user by their username.
//
// @Summary Get user by username
// @Description Retrieves user details based on their username.
// @Tags Users
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} response.StandardResponse "User data"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid username"
// @Failure 404 {object} response.StandardError "NOT FOUND - User not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router /users/get-user-by-name/{username} [get]
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

// CreateNewUser creates a new user in the system.
//
// @Summary Create new user
// @Description Creates a new user with the provided data.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body request.UserRequest true "New User Data"
// @Success 201 {object} string "CREATED - New user created successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 409 {object} response.StandardError "CONFLICT - User already exists"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Router /users/protected/create-new-user [post]
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

// UpdateUser updates an existing user's details.
//
// @Summary Update user
// @Description Updates an existing user's data by their ID.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body request.NewUser true "Updated User Data"
// @Success 200 {object} response.StandardResponse "UPDATED - User updated successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data"
// @Failure 404 {object} response.StandardError "NOT FOUND - User not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /users/protected/update-user/{id} [put]
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

// DeleteUser deletes a user by their ID.
//
// @Summary Delete user
// @Description Deletes a specific user by their unique identifier.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.StandardResponse "DELETED - User deleted successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid user ID"
// @Failure 404 {object} response.StandardError "NOT FOUND - User not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /users/protected/delete-user/{id} [delete]
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

// RestoreUser restores a previously deleted user by their ID.
//
// @Summary Restore user
// @Description Restores a previously deleted user by their unique identifier.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.StandardResponse "RESTORED - User restored successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid user ID"
// @Failure 404 {object} response.StandardError "NOT FOUND - User not found"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /users/protected/restore-user/{id} [patch]
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

// ChangeUserAvatar changes the avatar of a user by their ID.
//
// @Summary Change user avatar
// @Description Changes the avatar of a specific user by their ID.
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "User ID"
// @Param avatar formData file true "Avatar image file (PNG or JPEG, max 5MB)"
// @Success 200 {object} response.StandardResponse "Avatar uploaded successfully"
// @Failure 400 {object} response.StandardError "BAD REQUEST - Invalid input data or file type"
// @Failure 413 {object} response.StandardError "REQUEST ENTITY TOO LARGE - File size exceeds the limit"
// @Failure 415 {object} response.StandardError "UNSUPPORTED MEDIA TYPE - Only PNG and JPEG allowed"
// @Failure 500 {object} response.StandardError "INTERNAL SERVER ERROR - Unexpected server error"
// @Security AccessTokenAuth
// @Security RefreshTokenAuth
// @Router /users/protected/change-user-avatar/{id} [put]
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
