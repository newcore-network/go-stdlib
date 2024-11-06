
# Database Library Documentation

This document outlines the usage and setup of the `database` package within the `styerr-development/lib` repository. The library provides a unified way to establish database connections with support for different drivers like PostgreSQL, MySQL, and MariaDB.

## Overview

The `database` package simplifies the process of connecting to various databases using GORM as the ORM. The package is designed to be flexible, allowing easy integration and switching between database drivers.

## Features

- Support for PostgreSQL, MySQL, and MariaDB.
- Configurable connection using environment variables.
- Automatic database migration support.
- Easy-to-use interface for creating new connections.

## Installation

To use the `database` package in your project, import it as follows:

```go
import "github.com/styerr-development/lib/database"
```

## Usage

### 1. Configure Environment Variables

Ensure you have a `.env` file with the following variables:

```env
HOST=localhost
USER=gorm
PASSWORD=admin
DATABASE=database
PORT=5432
SSLMODE=false
JWTKEY="some text"
```

These configurations will be used to establish the database connection.

### 2. Create a Connection

Use the provided connection implementations to create a new database connection. Here is an example for connecting to PostgreSQL:

```go
package main

import (
    "log"
    "github.com/styerr-development/lib/database"
)

func main() {
    postgresConn := &database.PostgresConnection{}
    db, err := database.NewConnection(postgresConn, &YourModel{})
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }

    // Use `db` for your database operations
}
```

### 3. Switch to MySQL or MariaDB

To switch to MySQL or MariaDB, use the corresponding connection implementation:

**MySQL:**
```go
mysqlConn := &database.MySQLConnection{}
db, err := database.NewConnection(mysqlConn, &YourModel{})
```

**MariaDB:**
```go
mariaDBConn := &database.MariaDBConnection{}
db, err := database.NewConnection(mariaDBConn, &YourModel{})
```

## Connection Interface

The `database` package includes a `Connection` interface that must be implemented by any database type:

```go
type Connection interface {
    Connect(cfg configuration.GeneralConfig) (Conn, error)
}
```

## Auto Migration

To enable automatic migration of models, pass them as variadic arguments to the `NewConnection` function:

```go
db, err := database.NewConnection(postgresConn, &Model1{}, &Model2{})
```

## Error Handling

The `NewConnection` function attempts to connect to the database up to 5 times before failing and returning an error. Each failed attempt logs an error and waits for 3 seconds before retrying.
