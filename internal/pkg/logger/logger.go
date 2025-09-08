package logger

import (
	"os"

	"github.com/ray-d-song/go-echo-monolithic/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new logger based on configuration
func NewLogger(cfg *config.LoggerConfig) (*Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	var encoding zapcore.EncoderConfig
	if cfg.Encoding == "console" {
		encoding = zap.NewDevelopmentEncoderConfig()
	} else {
		encoding = zap.NewProductionEncoderConfig()
	}
	encoding.TimeKey = "timestamp"
	encoding.EncodeTime = zapcore.ISO8601TimeEncoder

	var output zapcore.WriteSyncer
	if cfg.OutputPath == "stdout" {
		output = zapcore.AddSync(os.Stdout)
	} else {
		file, err := os.OpenFile(cfg.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		output = zapcore.AddSync(file)
	}

	var encoder zapcore.Encoder
	if cfg.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoding)
	} else {
		encoder = zapcore.NewJSONEncoder(encoding)
	}

	core := zapcore.NewCore(encoder, output, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{Logger: logger}, nil
}