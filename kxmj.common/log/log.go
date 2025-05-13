package log

import (
	"go.uber.org/zap"
	"strings"
)

var (
	zLog *zap.Logger
)

type Config struct {
	Path            string `yaml:"path"`            //日志输出路径
	Prefix          string `yaml:"prefix"`          //日志文件前缀
	Level           string `yaml:"level"`           //日志级别：debug/info/error/warn
	Development     bool   `yaml:"development"`     //是否为开发者模式
	DebugFileSuffix string `yaml:"debugFileSuffix"` //debug日志文件后缀
	WarnFileSuffix  string `yaml:"warnFileSuffix"`  //warn日志文件后缀
	ErrorFileSuffix string `yaml:"errorFileSuffix"` //error日志文件后缀
	InfoFileSuffix  string `yaml:"infoFileSuffix"`  //info日志文件后缀
	MaxAge          int    `yaml:"maxAge"`          //保存的最大天数
	MaxBackups      int    `yaml:"maxBackups"`      //最多存在多少个切片文件
	MaxSize         int    `yaml:"maxSize"`         //日日志文件大小（M）
}

func create(config *Config) *zap.Logger {
	level := config.Level
	logLevel := zap.DebugLevel
	if strings.EqualFold("debug", level) {
		logLevel = zap.DebugLevel
	}

	if strings.EqualFold("info", level) {
		logLevel = zap.InfoLevel
	}

	if strings.EqualFold("error", level) {
		logLevel = zap.ErrorLevel
	}

	if strings.EqualFold("warn", level) {
		logLevel = zap.WarnLevel
	}

	return NewLogger(
		SetPath(config.Path),
		SetPrefix(config.Prefix),
		SetDevelopment(config.Development),
		SetDebugFileSuffix(config.DebugFileSuffix),
		SetWarnFileSuffix(config.WarnFileSuffix),
		SetErrorFileSuffix(config.ErrorFileSuffix),
		SetInfoFileSuffix(config.InfoFileSuffix),
		SetMaxAge(config.MaxAge),
		SetMaxBackups(config.MaxBackups),
		SetMaxSize(config.MaxSize),
		SetLevel(logLevel),
	)
}

func Init(config *Config) {
	zLog = create(config)
}

func Default() *zap.Logger {
	return zLog
}

func Sugar() *zap.SugaredLogger {
	return zLog.Sugar()
}
