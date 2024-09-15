package controller

import (
	"github.com/Dialosoft/src/domain/services"
	"github.com/gofiber/fiber/v3"
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
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "users not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}
