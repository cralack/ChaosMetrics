package xzap

import (
	"os"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Zap(env uint) (*zap.Logger, error) {
	if global.GVA_LOG != nil {
		return global.GVA_LOG, nil
	}
	// init val
	var (
		logConf zap.Config
		log     *zap.Logger
		err     error
	)
	// diff logger level
	switch env {
	case global.PRODUCT_ENV:
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
		getFileWriter(),
		logConf.Level,
	)

	// init logger
	switch env {
	case global.PRODUCT_ENV:
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
			// zap.AddStacktrace(logConf.Level),
		)
	}

	zap.ReplaceGlobals(log)

	log.Sync()

	if err != nil {
		panic(err)
	}

	return log, nil
}
