package log

import (
	ui "pgtoch/internal/UI"

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
	ui.PrintInfo(msg)
}

func (s *StyledLogger) Success(msg string, fields ...zap.Field) {
	s.logger.Info("Sucess", fields...)
	ui.PrintSuccess(msg)
}

func (s *StyledLogger) Error(msg string, fields ...zap.Field) {
	s.logger.Error(msg, fields...)
	ui.PrintError(msg)
}

func (s *StyledLogger) Warn(msg string, fields ...zap.Field) {
	s.logger.Warn(msg, fields...)
	ui.PrintWarning(msg)
}

func (s *StyledLogger) Fatal(msg string, fields ...zap.Field) {
	s.logger.Fatal(msg, fields...)
	ui.ExitWithError(msg)
}

func (s *StyledLogger) Debug(msg string, fields ...zap.Field) {
	s.logger.Debug(msg, fields...)
}

func (s *StyledLogger) Highligh(msg string, fields ...zap.Field) {
	s.logger.Info(msg, fields...)
	ui.PrintHighlight(msg)
}

func (s *StyledLogger) With(fields ...zap.Field) *StyledLogger {
	return &StyledLogger{
		logger: s.logger.With(fields...),
	}
}

func (s *StyledLogger) getZapLogger() *zap.Logger {
	return s.logger
}

var StyledLog *StyledLogger

func InitStyledLogger() {
	StyledLog = NewStyledLogger()
}
