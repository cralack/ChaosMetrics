package zlog

import (
	"fmt"
	"time"
	
	"github.com/cralack/ChaosMetrics/server/global"
	
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func getFileWriter() zapcore.WriteSyncer {
	logConf := global.GVA_CONF.LogConf
	logDir := global.GVA_CONF.DirTree.LogDir
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
