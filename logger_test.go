package stdlib

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitLoggerDevMode(t *testing.T) {
	InitLogger(true)
	assert.NotNil(t, logger, "Logger should be initialized in development mode")
}

func TestInitLoggerProdMode(t *testing.T) {
	InitLogger(false)
	assert.NotNil(t, logger, "Logger should be initialized in production mode")
}

func TestLoggingLevels(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)

	logger = zap.New(core)

	Info("Info message", map[string]interface{}{"key": "value"})
	Warn("Warn message", map[string]interface{}{"key": "value"})
	Error("Error message", map[string]interface{}{"key": "value"})

	output := buf.String()

	assert.Contains(t, output, "Info message", "Should contain'Info message'")
	assert.Contains(t, output, "Warn message", "Should contain 'Warn message'")
	assert.Contains(t, output, "Error message", "Should contain 'Error message'")
}

func TestCaptureError(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)

	logger = zap.New(core)

	CaptureError(simulatedError(), "Testing CaptureError", map[string]interface{}{
		"route": "/test",
	})

	output := buf.String()

	assert.Contains(t, output, "Testing CaptureError", "Should contain 'Testing CaptureError'")
	assert.Contains(t, output, "simulated error", "Should contain 'simulated error'")
}

func TestDebugLogging(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)

	logger = zap.New(core)

	Debug("Debug message", map[string]interface{}{"key": "debug_value"})

	output := buf.String()

	assert.Contains(t, output, "Debug message", "Should contain 'Debug message'")
	assert.Contains(t, output, "debug_value", "Should contain 'debug_value'")
}

type customError struct {
	Message string
}

func (e *customError) Error() string {
	return e.Message
}

// Helper: simulate an error
func simulatedError() error {
	return &customError{Message: "simulated error"}
}
