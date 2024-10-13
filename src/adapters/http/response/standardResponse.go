package response

import (
	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/gofiber/fiber/v3"
)

type StandardResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Standard(c fiber.Ctx, message string, data interface{}) error {
	std := StandardResponse{
		Message: message,
		Data:    data,
	}
	logger.Info(message, map[string]interface{}{
		"response": data,
		"message":  message,
		"route":    c.Path(),
		"method":   c.Method(),
	})
	return c.Status(fiber.StatusOK).JSON(std)
}

func StandardCreated(c fiber.Ctx, message string, data interface{}) error {
	std := StandardResponse{
		Message: message,
		Data:    data,
	}

	return c.Status(fiber.StatusCreated).JSON(std)
}
