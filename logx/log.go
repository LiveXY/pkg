package logx

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志配置
type LogConfig struct {
	Path  string `yaml:"path"`  // 日志路径
	Level string `yaml:"level"` // 日志记录级别
	Slow  int8   `yaml:"slow"`  // 慢查询 时长默认0秒 不记录慢查询日志
}

var (
	Logger *zap.Logger // 记录详细日志
	Error  *zap.Logger // 只记录错误日志
	DB     *zap.Logger // 只记录DB错误日志
)

// 默认初始化数据
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

// 异常处理
func Recover(desc string) {
	if err := recover(); err != nil {
		switch err.(type) {
		case string:
			Error.Error(desc + " 捕获到异常: " + err.(string))
		case error:
			Error.Error(desc + " 捕获到异常: ", zap.Error(errors.WithStack(err.(error))))
		}
	}
}

func Dump(vs ...any) {
	Logger.Info(spew.Sdump(vs))
}

// 记录错误
func Err(title string, err error) {
	Error.Error(title, zap.Error(errors.WithStack(err)))
}

func HasError(err error, titles ...string) {
	var title string
	if len(titles) > 0 {
		title = strings.Join(titles, " ")
	} else {
		title = "执行错误！"
	}
	if err != nil {
		Error.Error(title, zap.Error(errors.WithStack(err)))
	}
}

func getWriter(cfg LogConfig, name, ext string) lumberjack.Logger {
	today := time.Now().Format("20060102")
	logpath := path.Join(cfg.Path, name, today+"."+ext)
	return lumberjack.Logger{
		Filename:   logpath, // 日志文件路径
		MaxSize:    128,     // 每个日志文件保存的最大尺寸 单位：M  128
		MaxBackups: 30,      // 日志文件最多保存多少个备份 30
		MaxAge:     10,      // 文件最多保存多少天 7
		LocalTime:  true,    // 是否使用本地时间
		Compress:   false,   // 是否压缩
	}
}

func newInitLogger(cfg LogConfig, name, ext string) *zap.Logger {
	encodercfg := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		//CallerKey:     "line",
		MessageKey:    "msg",
		//StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	writer := getWriter(cfg, name, ext)
	consolecfg := zapcore.EncoderConfig{
		TimeKey:       "time",
		MessageKey:    "msg",
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
