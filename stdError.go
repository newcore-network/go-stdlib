package stdlib

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
)

type StandardError struct {
	ErrorMessage string `json:"error" example:"An error occurred - some context"`
}

type ValidatorError struct {
	ErrorMessage string            `json:"error"  example:"Validation failed"`
	Fields       map[string]string `json:"fields" example:"{'username': 'Username is required'}"`
}

// ErrInternalServer logs an internal server error and returns a JSON response
func ErrInternalServer(c fiber.Ctx, err error) error {
	response := StandardError{
		ErrorMessage: "INTERNAL SERVER ERROR",
	}
	CaptureError(err, "Internal server error occurred", map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusInternalServerError).JSON(response)
}

// ErrNotFound logs a not found error and returns a JSON response
func ErrNotFound(c fiber.Ctx) error {
	response := StandardError{
		ErrorMessage: "NOT FOUND",
	}
	Warn("Resource not found", map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusNotFound).JSON(response)
}

// ErrBadRequest logs a bad request error and returns a JSON response
func ErrBadRequest(c fiber.Ctx, err error) error {
	response := StandardError{
		ErrorMessage: "BAD REQUEST",
	}
	CaptureError(err, "Bad request error", map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusBadRequest).JSON(response)
}

// ErrConflict logs a conflict error and returns a JSON response
func ErrConflict(c fiber.Ctx, err error) error {
	response := StandardError{
		ErrorMessage: "CONFLICT",
	}
	Warn("Conflict error", map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusConflict).JSON(response)
}

// ErrUnauthorized logs an unauthorized error and returns a JSON response
func ErrUnauthorized(c fiber.Ctx, err error) error {
	response := StandardError{
		ErrorMessage: "UNAUTHORIZED",
	}
	Warn("Unauthorized access", map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusUnauthorized).JSON(response)
}

// ErrForbbiden logs a forbidden error and returns a JSON response
func ErrForbbiden(c fiber.Ctx, err error) error {
	response := StandardError{
		ErrorMessage: "FORBBIDEN",
	}
	Warn("forbidden access", map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusForbidden).JSON(response)
}

// ErrUUIDParse logs a bad UUID error and returns a JSON response
func ErrUUIDParse(c fiber.Ctx, id string) error {
	response := StandardError{
		ErrorMessage: "ID provided is not a valid UUID type",
	}
	Error(response.ErrorMessage, map[string]interface{}{
		"provided-id": id,
		"route":       c.Path(),
		"method":      c.Method(),
	})
	return c.Status(fiber.StatusBadRequest).JSON(response)
}

// ErrEmptyParametersOrArguments logs an error for missing parameters and returns a JSON response
func ErrEmptyParametersOrArguments(c fiber.Ctx) error {
	response := StandardError{
		ErrorMessage: "One of the parameters or arguments is empty",
	}
	Error(response.ErrorMessage, map[string]interface{}{
		"route":  c.Path(),
		"method": c.Method(),
	})
	return c.Status(fiber.StatusBadRequest).JSON(response)
}

// RegisterValidatorErr logs validation errors and returns a JSON response
func RegisterValidatorErr(c fiber.Ctx, errs error) error {
	response := ValidatorError{
		ErrorMessage: "Credential validation failed",
		Fields:       make(map[string]string),
	}
	validatorMessages := map[string]string{
		"Username": "Must be greater than 4 and less than 15",
		"Email":    "Invalid email",
		"Password": "Must be greater than 6",
	}
	for _, err := range errs.(validator.ValidationErrors) {
		response.Fields[err.Field()] = validatorMessages[err.Field()]
	}
	return c.Status(fiber.StatusBadRequest).JSON(response)
}

// PersonalizedErr returns an error with a custom message and status code
func PersonalizedErr(c fiber.Ctx, message string, status int) error {
	response := StandardError{
		ErrorMessage: message,
	}
	return c.Status(status).JSON(response)
}
