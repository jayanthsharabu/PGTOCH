package log

import "go.uber.org/zap"

var Logger *zap.Logger

func initLogger() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger:" + err.Error())
	}
	initStyledLogger()
}
