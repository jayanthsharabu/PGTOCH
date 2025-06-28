package log

import (
	"go.uber.org/zap"
)

type StyledLogger struct {
	logger *zap.Logger
}

func NewStyledLogger() *StyledLogger {
	return &StyledLogger{
		logger: Logger,
	}
}

func (s *StyledLogger) Info(msg string, fields ...zap.Field) {
	s.logger.Info(msg, fields...)

}
