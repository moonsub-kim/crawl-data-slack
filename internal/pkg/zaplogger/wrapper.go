package zaplogger

import (
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/logger"
	"go.uber.org/zap"
)

// Wrapper wraps zap logger to add context-dependent values into log
type Wrapper struct {
	logger Logger
}

// Info logs info-level log
func (w *Wrapper) Info(msg string, fields ...logger.Field) {
	zapFields := w.transformToZapFields(fields...)
	w.logger.Info(msg, zapFields...)
}

// Warn logs warn-level log
func (w *Wrapper) Warn(msg string, fields ...logger.Field) {
	zapFields := w.transformToZapFields(fields...)
	w.logger.Warn(msg, zapFields...)
}

// Error logs error-level log
func (w *Wrapper) Error(msg string, fields ...logger.Field) {
	zapFields := w.transformToZapFields(fields...)
	w.logger.Error(msg, zapFields...)
}

func (w *Wrapper) transformToZapFields(fields ...logger.Field) []zap.Field {
	zapField := make([]zap.Field, 0)
	for _, field := range fields {
		switch field.Value.(type) {
		case error:
			zapField = append(zapField, zap.Error(field.Value.(error)))
		default:
			zapField = append(zapField, zap.Any(field.Key, field.Value))
		}
	}
	return zapField
}

// NewZapLoggerWrapper creates Wrapper
func NewZapLoggerWrapper(logger Logger) *Wrapper {
	return &Wrapper{logger}
}
