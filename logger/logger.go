package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func (log *Logger) Fatal(err error) {
	log.Logger.Fatal(err.Error())
}

func New() (*Logger, error) {
	log, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("new logger: %w", err)
	}
	return &Logger{Logger: log}, nil
}
