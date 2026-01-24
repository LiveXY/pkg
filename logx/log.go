package logx

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig 日志库配置结构体
type LogConfig struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
	Slow  int8   `yaml:"slow"`
}

var (
	Logger = zap.NewNop()
	Error  = zap.NewNop()
	DB     = zap.NewNop()
)

// Init 初始化日志系统
func Init(cfg LogConfig, name ...string) {
	if len(name) > 0 {
		Logger = newInitLogger(cfg, name[0], "log")
		Error = newInitLogger(cfg, name[0], "err")
		DB = newInitLogger(cfg, name[0], "db")
	} else {
		Logger = newInitLogger(cfg, "", "log")
		Error = newInitLogger(cfg, "", "err")
		DB = newInitLogger(cfg, "", "db")
	}
}

// Recover 捕获并记录 Panic 异常
func Recover(desc string) {
	if err := recover(); err != nil {
		switch err := err.(type) {
		case string:
			Error.Error(desc + " 捕获到异常: " + err)
		case error:
			Error.Error(desc+" 捕获到异常: ", zap.Error(fmt.Errorf("%w", err)))
		}
	}
}

// Err 记录指定的错误信息
func Err(title string, err error) {
	Error.Error(title, zap.Error(fmt.Errorf("%w", err)))
}

// HasError 如果指定 err 不为 nil，则记录相关错误信息
func HasError(err error, titles ...string) {
	var title string
	if len(titles) > 0 {
		title = strings.Join(titles, " ")
	} else {
		title = "执行错误！"
	}
	if err != nil {
		Error.Error(title, zap.Error(fmt.Errorf("%w", err)))
	}
}

func getWriter(cfg LogConfig, name, ext string) lumberjack.Logger {
	today := time.Now().Format("20060102")
	logpath := path.Join(cfg.Path, name, today+"."+ext)
	return lumberjack.Logger{
		Filename:   logpath,
		MaxSize:    128,
		MaxBackups: 30,
		MaxAge:     10,
		LocalTime:  true,
		Compress:   false,
	}
}

func newInitLogger(cfg LogConfig, name, ext string) *zap.Logger {
	encodercfg := zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		NameKey:     "logger",
		MessageKey:  "msg",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	writer := getWriter(cfg, name, ext)
	consolecfg := zapcore.EncoderConfig{
		TimeKey:    "time",
		MessageKey: "msg",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
	}
	var level zapcore.Level
	switch cfg.Level {
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.DebugLevel
	}
	encoder := zapcore.NewConsoleEncoder(consolecfg)
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encodercfg),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(&writer)),
			zap.NewAtomicLevelAt(level),
		),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(zap.InfoLevel)),
	)
	caller := zap.AddCaller()
	development := zap.Development()
	return zap.New(core, caller, development)
}
