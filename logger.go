package stdlib

import (
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger(devMode bool) {
	var config zap.Config
	if devMode {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.MessageKey = "msg"
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	// Crear el logger
	var err error
	logger, err = config.Build()
	if err != nil {
		panic("Error init the logger system: " + err.Error())
	}
	defer logger.Sync()
}

// Info logs an informational message.
func Info(msg string, fields map[string]interface{}) {
	logger.With(createZapFields(fields)...).Info(msg)
}

// Warn logs a warning message.
func Warn(msg string, fields map[string]interface{}) {
	logger.With(createZapFields(fields)...).Warn(msg)
}

// Error logs an error message, including context such as file and line number.
func Error(msg string, fields map[string]interface{}) {
	addCallerInfo(fields)
	logger.With(createZapFields(fields)...).Error(msg)
}

// Debug logs a debug message, typically used for low-level system information.
func Debug(msg string, fields map[string]interface{}) {
	logger.With(createZapFields(fields)...).Debug(msg)
}

// Fatal logs a fatal error message and exits the application.
func Fatal(msg string, fields map[string]interface{}) {
	addCallerInfo(fields)
	logger.With(createZapFields(fields)...).Fatal(msg)
}

// CaptureError logs an error message with an additional error object.
func CaptureError(err error, msg string, fields map[string]interface{}) {
	if err != nil {
		fields["error"] = err.Error()
		Error(msg, fields)
	}
}

// Helper: Converts map[string]interface{} to []zap.Field
func createZapFields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

// Helper: Adds caller information to the fields
func addCallerInfo(fields map[string]interface{}) {
	if pc, file, line, ok := runtime.Caller(2); ok { // Saltar 2 niveles para capturar el caller correcto
		fn := runtime.FuncForPC(pc)
		fields["file"] = file
		fields["line"] = line
		fields["function"] = fn.Name()
	}
}
