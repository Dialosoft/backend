package response

import (
	"fmt"

	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/gofiber/fiber/v3"
)

type StandardError struct {
	ErrorMessage string `json:"error"`
}

func ErrInternalServer(c fiber.Ctx, err error, data interface{}, layer string) error {
	response := StandardError{
		ErrorMessage: "INTERNAL SERVER ERROR",
	}
	logger.CaptureError(err, fmt.Sprintf("(%s) Internal server error", layer), map[string]interface{}{
		"data":   data,
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusInternalServerError).JSON(response)
}

func ErrNotFound(c fiber.Ctx, layer string) error {
	err := StandardError{
		ErrorMessage: "NOT FOUND",
	}
	logger.Warn(fmt.Sprintf("(%s) Not found", layer), map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusNotFound).JSON(err)
}

func ErrBadRequest(c fiber.Ctx, body string, err error, layer string) error {
	response := StandardError{
		ErrorMessage: "BAD REQUEST",
	}
	logger.CaptureError(err, fmt.Sprintf("(%s) failed to parse | bad request!", layer), map[string]interface{}{
		"error":  err.Error(),
		"route":  c.Path(),
		"method": c.Method(),
		"body":   body,
	})
	return c.Status(fiber.StatusBadRequest).JSON(response)
}

func ErrBadRequestParse(c fiber.Ctx, err error, request interface{}, layer string) error {
	response := StandardError{
		ErrorMessage: "BAD REQUEST",
	}
	logger.CaptureError(err, fmt.Sprintf("(%s) Failed to parse %v", layer, request), map[string]interface{}{
		"request-tried": request,
		"route":         c.Path(),
		"method":        c.Method(),
	})
	return c.Status(fiber.StatusBadRequest).JSON(response)
}

func ErrConflict(c fiber.Ctx, err error, request interface{}, layer string) error {
	response := StandardError{
		ErrorMessage: "CONFLICT",
	}
	logger.Warn(fmt.Sprintf("(%s) Conflict! error: %v", layer, err), map[string]interface{}{
		"request-tried": request,
		"route":         c.Path(),
		"method":        c.Method(),
	})
	return c.Status(fiber.StatusConflict).JSON(response)
}

func ErrUnauthorized(c fiber.Ctx, data interface{}, err error, layer string) error {
	response := StandardError{
		ErrorMessage: "UNAUTHORIZED",
	}
	logger.Warn(fmt.Sprintf("(%s) Unauthorized!, error: %v", layer, err), map[string]interface{}{
		"data":   data,
		"error":  err.Error(),
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusUnauthorized).JSON(response)
}

func ErrExpiredAccessToken(c fiber.Ctx) error {
	err := StandardError{
		ErrorMessage: "AccessToken expired",
	}
	return c.Status(fiber.StatusUnauthorized).JSON(err)
}

func ErrForbidden(c fiber.Ctx) error {
	err := StandardError{
		ErrorMessage: "FORBIDDEN",
	}
	return c.Status(fiber.StatusForbidden).JSON(err)
}

func ErrUnauthorizedHeader(c fiber.Ctx) error {
	err := StandardError{
		ErrorMessage: "Authorization Header is missing",
	}
	return c.Status(fiber.StatusUnauthorized).JSON(err)
}

func ErrUnauthorizedInvalidHeader(c fiber.Ctx) error {
	err := StandardError{
		ErrorMessage: "Invalid authorization header format",
	}
	return c.Status(fiber.StatusUnauthorized).JSON(err)
}

func ErrUUIDParse(c fiber.Ctx, id string) error {
	response := StandardError{
		ErrorMessage: "ID provided is not a valid UUID type",
	}
	logger.Error(response.ErrorMessage, map[string]interface{}{
		"provided-id": id,
		"route":       c.Path(),
		"method":      c.Method(),
	})
	return c.Status(fiber.StatusBadRequest).JSON(response)
}

func ErrEmptyParametersOrArguments(c fiber.Ctx) error {
	response := StandardError{
		ErrorMessage: "One of the parameters or arguments is empty",
	}
	logger.Error(response.ErrorMessage, map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusBadRequest).JSON(response)
}

func PersonalizedErr(c fiber.Ctx, message string, status int) error {
	err := StandardError{
		ErrorMessage: message,
	}
	return c.Status(status).JSON(err)
}
