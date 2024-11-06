
# Logger Package Documentation

This document provides an overview and usage guide for the `logger` package within your Go project. The `logger` package is built using `logrus` and `lumberjack` for structured and robust logging capabilities, including file rotation and JSON formatting.

## Overview

The `logger` package offers a centralized logging utility for your application with the following features:
- JSON-formatted logs.
- Multi-output logging (console and file).
- Log rotation with customizable settings.
- Captures caller information (file, line, and function) for error and fatal logs.
- Supports different log levels: `Info`, `Warn`, `Error`, `Debug`, and `Fatal`.

## Initialization

To start using the `logger` package, you need to initialize it by calling `InitLogger()` in your `main()` function or during the application's startup sequence:

```go
package main

import (
    "your_project/logger"
)

func main() {
    logger.InitLogger()
    // Application code...
}
```

### Log File Configuration

The package writes logs to both the console and a rotating log file (`./logs/app.log`). The file rotation is configured as follows:
- **MaxSize**: 50 MB per log file.
- **MaxBackups**: Retains up to 3 backup log files.
- **MaxAge**: Retains log files for 28 days.
- **Compress**: Compresses old log files.

## Usage

### Logging Methods

- **Info**: Logs an informational message.
  ```go
  logger.Info("Application started", map[string]interface{}{
      "version": "1.0.0",
  })
  ```

- **Warn**: Logs a warning message.
  ```go
  logger.Warn("Memory usage is high", map[string]interface{}{
      "threshold": "80%",
  })
  ```

- **Error**: Logs an error message, including the file and line number.
  ```go
  logger.Error("Failed to load configuration", map[string]interface{}{
      "configPath": "/etc/app/config.yaml",
  })
  ```

- **Debug**: Logs a debug message for low-level details.
  ```go
  logger.Debug("Debugging connection", map[string]interface{}{
      "host": "localhost",
      "port": 5432,
  })
  ```

- **Fatal**: Logs a fatal error and exits the application.
  ```go
  logger.Fatal("Unrecoverable error occurred", map[string]interface{}{
      "reason": "Out of memory",
  })
  ```

- **CaptureError**: Logs an error message along with an `error` object.
  ```go
  err := errors.New("unexpected failure")
  logger.CaptureError(err, "Error processing request", map[string]interface{}{
      "requestID": "12345",
  })
  ```

## Capturing Caller Information

The `Error` and `Fatal` methods automatically capture the file, line, and function name where the log was called. This helps in tracing the exact location of the issue in the codebase.

## Dependencies

The `logger` package relies on the following third-party libraries:
- `github.com/sirupsen/logrus`
- `gopkg.in/natefinch/lumberjack.v2`

Ensure these dependencies are included in your `go.mod` file.
