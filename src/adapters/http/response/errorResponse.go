package response

import (
	"fmt"

	"github.com/Dialosoft/src/pkg/utils/logger"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
)

type StandardError struct {
	ErrorMessage string `json:"error"`
}

type ValidatorError struct {
	ErrorMessage string            `json:"error"`
	Fields       map[string]string `json:"fields"`
}

// ErrInternalServer returns an error with the message "INTERNAL SERVER ERROR"
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

// ErrNotFound returns an error with the message "NOT FOUND"
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

// ErrBadRequest returns an error with the message "BAD REQUEST"
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

// ErrBadRequestParse returns an error with the message "BAD REQUEST"
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

// ErrConflict returns an error with the message "CONFLICT"
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

// ErrUnauthorized returns an error with the message "UNAUTHORIZED"
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

// ErrExpiredAccessToken returns an error with the message "AccessToken expired"
func ErrExpiredAccessToken(c fiber.Ctx, layer string) error {
	err := StandardError{
		ErrorMessage: "AccessToken expired",
	}
	logger.Warn(fmt.Sprintf("(%s) Unauthorized | AccessToken expired", layer), map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusUnauthorized).JSON(err)
}

// ErrInvalidToken returns an error with the message "Invalid token"
func ErrInvalidToken(c fiber.Ctx, layer string) error {
	err := StandardError{
		ErrorMessage: "Invalid token",
	}
	logger.Warn(fmt.Sprintf("(%s) Unauthorized | Invalid token", layer), map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusUnauthorized).JSON(err)
}

// ErrTokenIsBlacklisted returns an error with the message "Token is blacklisted"
func ErrTokenIsBlacklisted(c fiber.Ctx, layer string) error {
	response := StandardError{
		ErrorMessage: "Token is blacklisted",
	}
	logger.Warn(fmt.Sprintf("(%s) Unauthorized | Token is blacklisted", layer), map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusUnauthorized).JSON(response)
}

// ErrForbidden returns an error with the message "FORBIDDEN"
func ErrForbidden(c fiber.Ctx, layer string) error {
	err := StandardError{
		ErrorMessage: "FORBIDDEN",
	}
	logger.Warn(fmt.Sprintf("(%s) Forbidden | insufficient permissions", layer), map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusForbidden).JSON(err)
}

// ErrUnauthorizedHeader returns an error with the message "Authorization Header is missing"
func ErrUnauthorizedHeader(c fiber.Ctx, layer string) error {
	err := StandardError{
		ErrorMessage: "Authorization Header is missing",
	}
	logger.Warn(fmt.Sprintf("(%s) Unauthorized | Authorization Header is missing", layer), map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusUnauthorized).JSON(err)
}

// ErrUnauthorizedInvalidHeader returns an error with the message "Invalid authorization header format"
func ErrUnauthorizedInvalidHeader(c fiber.Ctx, layer string) error {
	err := StandardError{
		ErrorMessage: "Invalid authorization header format",
	}
	logger.Warn(fmt.Sprintf("(%s) Unauthorized | Invalid authorization header format", layer), map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusUnauthorized).JSON(err)
}

// ErrUUIDParse returns an error with the message "ID provided is not a valid UUID type"
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

// ErrEmptyParametersOrArguments returns an error with the message "One of the parameters or arguments is empty"
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

// RegisterValidatiorErr returns an error with the message "Credential validation failed"
func RegisterValidatiorErr(c fiber.Ctx, errs error) error {
	err := ValidatorError{
		ErrorMessage: "Credential validation failed",
		Fields:       make(map[string]string),
	}
	validatorMesssages := map[string]string{
		"Username": "Must be greater than 4 and less than 15",
		"Email":    "Invalid email",
		"Password": "Must be greater than 6",
	}
	for _, error := range errs.(validator.ValidationErrors) {
		err.Fields[error.Field()] = validatorMesssages[error.Field()]
	}
	return c.Status(fiber.StatusBadRequest).JSON(err)
}

// PersonalizedErr returns an error with the message "message" and the status "status"
func PersonalizedErr(c fiber.Ctx, message string, status int) error {
	err := StandardError{
		ErrorMessage: message,
	}
	return c.Status(status).JSON(err)
}
