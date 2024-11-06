
# Standard Responses Package Documentation

This document outlines the `standardResponses` package, which provides standardized methods for generating consistent API responses and error handling in applications using the Fiber web framework.

## Overview

The `standardResponses` package ensures uniformity in API responses by offering a set of pre-defined response and error handling functions. These functions help streamline logging and response formatting, while also including relevant context about the route, method, and data involved.

## Features

- Pre-defined methods for common HTTP responses such as `OK`, `Created`, `BadRequest`, `Unauthorized`, etc.
- Structured error and success responses with JSON formatting.
- Integration with the `logger` package for automatic logging of responses and errors.
- Validation error handling and field-specific messaging support.

## Response Structs

### `StandardError`

Used for error responses:
```go
type StandardError struct {
    ErrorMessage string `json:"error" example:"An error occurred - some context"`
}
```

### `ValidatorError`

Used for validation error responses:
```go
type ValidatorError struct {
    ErrorMessage string            `json:"error" example:"Validation failed"`
    Fields       map[string]string `json:"fields" example:"{'username': 'Username is required'}"`
}
```

### `StandardResponse`

Used for standard success responses:
```go
type StandardResponse struct {
    Message string      `json:"message" example:"info message"`
    Data    interface{} `json:"data"`
}
```

## Error Handling Methods

### `ErrInternalServer`

Returns an internal server error response:
```go
func ErrInternalServer(c fiber.Ctx, err error, data interface{}, layer string) error
```

### `ErrNotFound`

Returns a not found error response:
```go
func ErrNotFound(c fiber.Ctx, layer string) error
```

### `ErrBadRequest`

Returns a bad request error response:
```go
func ErrBadRequest(c fiber.Ctx, body string, err error, layer string) error
```

### `ErrConflict`

Returns a conflict error response:
```go
func ErrConflict(c fiber.Ctx, err error, request interface{}, layer string) error
```

### `ErrUnauthorized`

Returns an unauthorized error response:
```go
func ErrUnauthorized(c fiber.Ctx, data interface{}, err error, layer string) error
```

### Other Error Methods

The package includes additional error methods such as `ErrExpiredAccessToken`, `ErrInvalidToken`, `ErrTokenIsBlacklisted`, `ErrForbidden`, `ErrUnauthorizedHeader`, and more.

## Success Response Methods

### `Standard`

Returns a standard `200 OK` response with a message and data:
```go
func Standard(c fiber.Ctx, message string, data interface{}) error
```

### `StandardCreated`

Returns a `201 Created` response with a message and data:
```go
func StandardCreated(c fiber.Ctx, message string, data interface{}) error
```

## Example Usage

### Error Response Example
```go
if err := someFunction(); err != nil {
    return standardResponses.ErrInternalServer(c, err, requestData, "HandlerName")
}
```

### Success Response Example
```go
return standardResponses.Standard(c, "Data retrieved successfully", data)
```

## Logging Integration

Each response method logs relevant context automatically using the `logger` package. For example, `ErrInternalServer` logs the error with context about the route, method, and additional data provided.

## Dependencies

The `standardResponses` package relies on:
- `github.com/gofiber/fiber/v3`: The Fiber web framework.
- `github.com/go-playground/validator`: For handling validation errors.
- `github.com/styerr-development/libs/logger`: For integrated logging.
