package main

import (
	"github.com/yuanyp8/log"
)

func defaultTeeLogger() {
	log.ResetDefault(log.DefaultTee())
	defer log.Sync()

	for i := 0; i < 20000; i++ {
		log.Info("rotate_logger:",
			log.String("app", "start ok"),
			log.Int("major version", 2),
		)
		log.Error("rotate_logger:",
			log.String("app", "crash"),
			log.Int("reason", -1),
		)
	}
}

func customTeeLogger() {
	tops := []log.TeeOption{
		{Filename: "custom/access.log",
			Ropt: log.RotateOptions{
				MaxSize:    2,
				MaxAge:     2,
				MaxBackups: 2,
				Compress:   false,
			},
			Lef: func(lvl log.Level) bool {
				return lvl <= log.InfoLevel
			},
		},
		{Filename: "custom/error.log",
			Ropt: log.RotateOptions{
				MaxSize:    2,
				MaxAge:     2,
				MaxBackups: 2,
				Compress:   false,
			},
			Lef: func(lvl log.Level) bool {
				return lvl > log.InfoLevel
			},
		},
	}

	logger := log.NewLoggerTeeWithRotate(tops)
	log.ResetDefault(logger)
	defer log.Sync()

	for i := 0; i < 20000; i++ {
		log.Info("rotate_logger:",
			log.String("app", "start ok"),
			log.Int("major version", 2),
		)
		log.Error("rotate_logger:",
			log.String("app", "crash"),
			log.Int("reason", -1),
		)
	}
}

func main() {
	defaultTeeLogger()
	customTeeLogger()
}
