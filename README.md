# Styerr Standard Libraries

Welcome to the Styerr Standard Library, a comprehensive suite of reusable components designed for seamless integration across all Styerr repositories and services. These libraries offer robust functionality, streamline development, and promote consistency throughout projects.

## Overview
The Styerr Standard Library includes key packages for common functionalities such as configuration management, database connections, structured logging, and standardized HTTP responses. Each package is designed to be modular and easy to integrate into different projects.

## Available libraries
Explore the available libraries below, each with its own documentation and usage examples:
- [configuration](https://github.com/styerr-development/libs/blob/master/configuration/configuration.md): Manage and load application configuration from environment files.
- [database](https://github.com/styerr-development/libs/blob/master/database/database.md): Establish and manage database connections with support for PostgreSQL, MySQL, and MariaDB.
- [logger](https://github.com/styerr-development/libs/blob/master/logger/logger.md): A structured logging package with support for JSON formatting and log file rotation.
- [standardResponses](https://github.com/styerr-development/libs/blob/master/standardResponses/standardResponses.md): Simplifies the creation of standardized HTTP responses and error handling for APIs.

## Getting Started
### Prerequisites
Ensure you have the following in place before using the libraries:
- .env file: Required for loading configuration details such as database credentials. Refer to the .env.example file in the repository for a template to get started.
### Installation
Add the desired libraries to your project using Go modules:

```sh
$ go get github.com/styerr-development/libs
```
for all libraries



## Simple configuration
To set up your application using the configuration library:

- Create a .env file: This file should contain all necessary environment variables for your application, such as database host, user, password, etc.
- Load configuration: Use the configuration package to load and manage your application settings.

```conf
DB_HOST=localhost
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=mydatabase
DB_PORT=5432
DB_SSLMODE=disable
JWT_KEY=myjwtsecretkey
```

### Usage in your Go application:

```go
import (
    "github.com/styerr-development/libs/configuration"
)

func main() {
    config := configuration.GetFromEnvFile(".env")
    // Now you can use the `config` object in your application
    // -> Check the list of libraries to review their documentation and learn how to implement each one... <-
}
```
