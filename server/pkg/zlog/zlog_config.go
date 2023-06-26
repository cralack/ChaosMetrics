package zlog

import (
	"ChaosMetrics/server/global"
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DebugLevel  = zapcore.DebugLevel
	InfoLevel   = zapcore.InfoLevel
	WarnLevel   = zapcore.WarnLevel
	ErrorLevel  = zapcore.ErrorLevel
	DPanicLevel = zapcore.DPanicLevel
	PanicLevel  = zapcore.PanicLevel
	FatalLevel  = zapcore.FatalLevel
)

type LoggerConfig struct {
	MaxSize    int  `mapstructure:"maxsize" yaml:"maxsize"`       // 日志的最大大小，以M为单位
	MaxBackups int  `mapstructure:"maxbackups" yaml:"maxbackups"` // 保留的旧日志文件的最大数量
	MaxAge     int  `mapstructure:"maxage" yaml:"maxage"`         // 保留旧日志文件的最大天数
	LocalTime  bool `mapstructure:"localtime" yaml:"localtime"`   // 是否使用本地时间
	Compress   bool `mapstructure:"compress" yaml:"compress"`     // 是否压缩旧日志文件
}

func getFileWriter() zapcore.WriteSyncer {
	// var logConf *config.LoggerConfig
	// var err error
	// //parse to struct
	// err = global.GVA_VP.UnmarshalKey("logger", &logConf)
	// if err != nil {
	// 	global.GVA_LOG.Fatal("viper unmarshal db config failed")
	// }
	// global.GVA_CONF.LogConf = logConf
	logConf := global.GVA_CONF.LogConf
	logDir := global.GVA_CONF.DirTree.LogDIr
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
