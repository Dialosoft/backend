package response

import "github.com/gofiber/fiber/v3"

type StandardError struct {
	ErrorMessage string `json:"error"`
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
