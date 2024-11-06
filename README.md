# Database: Styerr Standard Library

A simple Go library for establishing a database connection using `gorm`, designed to standardize connection setups. Currently supports PostgreSQL and includes retry logic.


## Installation

```sh
$   go get -u github.com/styerr-development/library-database
```

## Usage

```go
package main

import (
	"log"
	"github.com/styerr-development/library-database"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	
	// `db` is ready for use
}

```

## Configuration

Requires a .env file with database credentials. The configuration is handled through library-configuration. Check .env.example file for example

## Notes

- Supports retrying connection up to five times with a 3-second delay.
- Planned support for additional databases in the future.
