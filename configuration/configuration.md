
# Configuration Package Documentation

This document outlines the `configuration` package, which provides functionality for reading and managing application configuration from environment files.

## Overview

The `configuration` package simplifies the process of loading environment variables and converting them into a structured `GeneralConfig` type. This ensures that your application has a consistent and centralized way to manage configuration settings.

## Features

- Loads configuration from a `.env` file.
- Provides structured configuration with the `GeneralConfig` type.
- Handles environment variable parsing, including integer and conditional checks.
- Logs meaningful errors when required environment variables are missing or misconfigured.

## Dependencies

The `configuration` package relies on:
- `github.com/joho/godotenv`: For loading `.env` files into the environment.

Ensure this dependency is included in your `go.mod` file.

## GeneralConfig Struct

The `GeneralConfig` struct holds the main configuration properties:

```go
type GeneralConfig struct {
    Host     string
    User     string
    Password string
    Database string
    Port     int
    SSLMode  string
    JWTKey   string
}
```

## Functionality

### `GetFromEnvFile`

The `GetFromEnvFile` function loads environment variables from the specified `.env` file and returns a `GeneralConfig` instance. If any critical variable is missing or misconfigured, it logs an error and stops the execution.

**Function signature:**
```go
func GetFromEnvFile(file string) GeneralConfig
```

### Example Usage

Here's how you can use the `GetFromEnvFile` function in your application:

```go
package main

import (
    "log"
    "github.com/newcore-network/libsconfiguration"
)

func main() {
    config := configuration.GetFromEnvFile(".env")
    log.Printf("Configuration loaded: %+v\n", config)
}
```

### Environment Variables

The package expects the following environment variables:

- **HOST**: The database host address.
- **USER**: The username for the database connection.
- **PASSWORD**: The password for the database connection.
- **DATABASE**: The name of the database.
- **PORT**: The port number for the database connection (as an integer).
- **SSLMODE**: SSL mode for the database (e.g., "enable" or "disable").
- **JWTKEY**: The secret key for JWT authentication.

### Error Handling

- The function logs an error and exits if the `.env` file cannot be loaded.
- Logs a warning if the `PORT` variable is invalid and returns an empty `GeneralConfig`.
- Ensures `JWTKEY` is present; otherwise, it logs a fatal error and stops execution.

## Constants

The `configuration` package uses constants from `constants` for standardized key names and error messages.
