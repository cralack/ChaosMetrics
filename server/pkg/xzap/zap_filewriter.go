package xzap

import (
	"fmt"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func newFileWriter() zapcore.WriteSyncer {
	logConf := global.ChaConf.LogConf
	logDir := global.ChaConf.DirTree.LogDir
	today := time.Now().Format("20060102") // 格式化为年月日，如20230626
	filename := fmt.Sprintf("%s/log_%s.log", logDir, today)
	writer := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   filename,
			MaxSize:    logConf.MaxSize,
			MaxBackups: logConf.MaxBackups,
			MaxAge:     logConf.MaxAge,
			LocalTime:  logConf.LocalTime,
			Compress:   logConf.Compress,
		})
	return writer
}
