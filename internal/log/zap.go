package log

import (
	"discord-music/internal/config"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = "debgu"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

type ZapLogger struct {
	cfg    *config.Config
	logger *zap.Logger
}

// Debug implements Logger.
func (z *ZapLogger) Debug(args ...any) {
	z.logger.Sugar().Debug(args)
}

// Error implements Logger.
func (z *ZapLogger) Error(args ...any) {
	z.logger.Sugar().Error(args)
}

// Fatal implements Logger.
func (z *ZapLogger) Fatal(args ...any) {
	z.logger.Sugar().Fatal(args)
}

// Info implements Logger.
func (z *ZapLogger) Info(args ...any) {
	z.logger.Sugar().Info(args)
}

// Warn implements Logger.
func (z *ZapLogger) Warn(args ...any) {
	z.logger.Sugar().Warn(args)
}

// Fatalf implements Logger.
func (z *ZapLogger) Fatalf(pattern string, args ...any) {
	z.logger.Sugar().Fatalf(pattern, args)
}

// Debugf implements Logger.
func (z *ZapLogger) Debugf(pattern string, args ...any) {
	z.logger.Sugar().Debugf(pattern, args)
}

// Errorf implements Logger.
func (z *ZapLogger) Errorf(pattern string, args ...any) {
	z.logger.Sugar().Errorf(pattern, args)
}

// Infof implements Logger.
func (z *ZapLogger) Infof(pattern string, args ...any) {
	z.logger.Sugar().Infof(pattern, args)
}

// Warnf implements Logger.
func (z *ZapLogger) Warnf(pattern string, args ...any) {
	z.logger.Sugar().Warnf(pattern, args)
}

func convertLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	default:
		fallthrough
	case DebugLevel:
		return zap.DebugLevel
	case InfoLevel:
		return zap.InfoLevel
	case WarnLevel:
		return zap.WarnLevel
	case ErrorLevel:
		return zap.ErrorLevel
	}
}

func NewZapLogger(cfg *config.Config) (Logger, error) {
	// for docker
	logfilePath := "/var/app.log"

	outputPaths := []string{"stdout"}
	if !cfg.DEV {
		outputPaths = append(outputPaths, logfilePath)
	}

	level := convertLogLevel(cfg.LOG_LEVEL)

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      cfg.DEV,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      outputPaths,
		ErrorOutputPaths: outputPaths,
	}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		logger: logger,
		cfg:    cfg,
	}, nil
}
