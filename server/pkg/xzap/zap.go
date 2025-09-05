package xzap

import (
	"os"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Zap(env global.Env) (*zap.Logger, error) {
	if global.ChaLogger != nil {
		return global.ChaLogger, nil
	}
	// init val
	var (
		logConf zap.Config
		log     *zap.Logger
		err     error
	)
	// diff logger level
	switch env {
	case global.ProductEnv:
		logConf = zap.NewProductionConfig()
	default:
		logConf = zap.NewDevelopmentConfig()
	}

	// setup logConf
	// logConf.EncoderConfig.TimeKey = "timestamp"
	logConf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	// init log console core
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(logConf.EncoderConfig),
		zapcore.Lock(os.Stdout),
		logConf.Level,
	)

	// init log file core
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(logConf.EncoderConfig),
		newFileWriter(),
		logConf.Level,
	)

	// init logger
	switch env {
	case global.ProductEnv:
		log = zap.New(
			zapcore.NewTee(
				consoleCore,
				fileCore,
			),
		)
	default:
		log = zap.New(
			zapcore.NewTee(
				consoleCore,
			),
			zap.AddCaller(),
		)
	}

	zap.ReplaceGlobals(log)

	defer func(log *zap.Logger) {
		err = log.Sync()
	}(log)

	if err != nil {
		panic(err)
	}

	return log, nil
}
