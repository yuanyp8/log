package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"sync"
	"time"
)

var (
	std       *Logger = NewDefaultLogger()
	teeLogger *Logger = NewLoggerTeeWithRotate(defaultTops)
	once              = sync.Once{}
)

type LevelEnablerFunc func(lvl Level) bool

type Logger struct {
	l     *zap.Logger // zap ensure that zap.Logger is safe for concurrent use
	level Level
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func Sync() (err error) {
	if std != nil {
		once.Do(func() {
			err = std.Sync()
		})
		return
	}
	return
}

// NewDefaultLogger 对zap的封装, 将创建 `Logger` 的过程封装到 `New` 内部
func NewDefaultLogger() *Logger {
	return New(os.Stderr, InfoLevel, WithCaller(true))
}

// ResetDefault not safe for concurrent use
func ResetDefault(l *Logger) {
	std = l
	Info = std.Info
	Warn = std.Warn
	Error = std.Error
	DPanic = std.DPanic
	Panic = std.Panic
	Fatal = std.Fatal
	Debug = std.Debug
}

// New create a new logger (not support log rotating).
func New(writer io.Writer, level Level, opts ...Option) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(timeFormat))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(writer),
		zapcore.Level(level),
	)
	logger := &Logger{
		l:     zap.New(core, opts...),
		level: level,
	}
	return logger
}

func Default() *Logger {
	return std
}

func DefaultTee() *Logger {
	return teeLogger
}
