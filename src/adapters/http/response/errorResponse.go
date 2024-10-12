package response

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type StandardError struct {
	ErrorMessage string `json:"error"`
}

type ValidatorError struct {
	ErrorMessage string            `json:"error"`
	Fields       map[string]string `json:"fields"`
}

func ErrInternalServer(c fiber.Ctx) error {
	err := StandardError{
		ErrorMessage: "INTERNAL SERVER ERROR",
	}
	return c.Status(fiber.StatusInternalServerError).JSON(err)
}

func ErrNotFound(c fiber.Ctx) error {
	err := StandardError{
		ErrorMessage: "NOT FOUND",
	}
	return c.Status(fiber.StatusNotFound).JSON(err)
}

func ErrBadRequest(c fiber.Ctx) error {
	err := StandardError{
		ErrorMessage: "BAD REQUEST",
	}
	return c.Status(fiber.StatusBadRequest).JSON(err)
}

func ErrConflict(c fiber.Ctx) error {
	err := StandardError{
		ErrorMessage: "CONFLICT",
	}
	return c.Status(fiber.StatusConflict).JSON(err)
}

func ErrUnauthorized(c fiber.Ctx) error {
	err := StandardError{
		ErrorMessage: "UNAUTHORIZED",
	}
	return c.Status(fiber.StatusUnauthorized).JSON(err)
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

func ErrUUIDParse(c fiber.Ctx) error {
	err := StandardError{
		ErrorMessage: "ID provided is not a valid UUID type",
	}
	return c.Status(fiber.StatusBadRequest).JSON(err)
}

func ErrEmptyParametersOrArguments(c fiber.Ctx) error {
	err := StandardError{
		ErrorMessage: "One of the parameters or arguments is empty",
	}
	return c.Status(fiber.StatusBadRequest).JSON(err)
}

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

func PersonalizedErr(c fiber.Ctx, message string, status int) error {
	err := StandardError{
		ErrorMessage: message,
	}
	return c.Status(status).JSON(err)
}
