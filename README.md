# Newcore.gg Standard Libraries

Welcome to the Newcore Standard Library, a comprehensive suite of reusable components designed for seamless integration across all Newcore repositories and services. These libraries offer robust functionality, streamline development, and promote consistency throughout projects.

Test Coverage: 11.3%

## Overview
The Newcore Standard Library includes key packages for common functionalities such as configuration management, database connections, structured logging, and standardized HTTP responses. Each package is designed to be modular and easy to integrate into different projects.

## Docs
Explore the available utilities in `docs/` folder
- [configuration](https://github.com/newcore-network/go-stdlib/docs/configuration.md): Manage and load application configuration from environment files.
- [database](https://github.com/newcore-network/go-stdlib/docs/database.md): Establish and manage database connections with support for PostgreSQL, MySQL, and MariaDB.
- [logger](https://github.com/newcore-network/go-stdlib/docs/logger.md): A structured logging package with support for JSON formatting and log file rotation.
- [stdResponses](https://github.com/newcore-network/go-stdlib/docs/stdResponses.md): Simplifies the creation of standardized HTTP responses and error handling for APIs.

## Getting Started
### Prerequisites
Ensure you have the following in place before using the libraries:
- Using Fiber version 3.x
- Go version 1.23 or higher: Ensure you have the latest version of Go installed on your system.
- .env file: Required for loading configuration details such as database credentials. Refer to the .env.example file in the repository for a template to get started.
### Installation
Add the desired libraries to your project using Go modules:

```sh
$ go get github.com/newcore-network/go-stdlib
```
for all libraries



## Simple configuration
To set up your application using the configuration library:

- Create a .env file: This file should contain all necessary environment variables for your application, such as database host, user, password, etc.
- Load configuration: Use the configuration package to load and manage your application settings.

```conf
#Configure the database connection
POSTGRES_HOST=postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=admin
POSTGRES_DATABASE=altv
POSTGRES_PORT=5432
POSTGRES_SSLMODE=false

# Configure the redis connection

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=admin
REDIS_DB=0
```

### Usage in your Go application:

```go
import (
    "github.com/newcore-network/go-stdlib"
)

func main() {
    config := stdlib.LoadCfg(".env")
    // Now you can use the `config` object in your application
    // -> Check the list of libraries to review their documentation and learn how to implement each one... <-
}
```
