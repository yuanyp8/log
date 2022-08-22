package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

type Option = zap.Option

var (
	WithCaller    = zap.WithCaller
	AddStacktrace = zap.AddStacktrace
)

// set default TeeOption list called tops;
// access logs will be writing into access.log
// error logs will be writing into error.log
// the files max size is 1 Mb; max age is 1 day;
// the files max rotate replicaset is 3;
// and the log files will be composing with tar.gz format
var defaultTops = []TeeOption{
	{Filename: "access.log",
		Ropt: RotateOptions{
			MaxSize:    1,
			MaxAge:     1,
			MaxBackups: 3,
			Compress:   true,
		},
		Lef: func(lvl Level) bool {
			return lvl <= InfoLevel
		},
	},
	{Filename: "error.log",
		Ropt: RotateOptions{
			MaxSize:    1,
			MaxAge:     1,
			MaxBackups: 1,
			Compress:   true,
		},
		Lef: func(lvl Level) bool {
			return lvl > InfoLevel
		},
	},
}

type TeeOption struct {
	Filename string
	Ropt     RotateOptions
	Lef      LevelEnablerFunc
}

type RotateOptions struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

func NewLoggerTeeWithRotate(tops []TeeOption, opts ...Option) *Logger {
	var cores []zapcore.Core
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(timeFormat))
	}

	for _, top := range tops {
		top := top

		lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return top.Lef(lvl)
		})

		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   top.Filename,
			MaxSize:    top.Ropt.MaxSize,
			MaxBackups: top.Ropt.MaxBackups,
			MaxAge:     top.Ropt.MaxAge,
			Compress:   top.Ropt.Compress,
		})

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.AddSync(w),
			lv,
		)
		cores = append(cores, core)
	}

	logger := &Logger{
		l: zap.New(zapcore.NewTee(cores...), opts...),
	}
	return logger
}
