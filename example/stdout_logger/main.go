package main

import (
	"fmt"
	"github.com/yuanyp8/log"
	"os"
)

func defaultLogger() {
	defer log.Sync()

	log.Info("default stdout logger",
		log.String("level", "info"),
		log.String("app", "start ok"),
	)

	log.Debug("custom stdout logger",
		log.String("level", "debug"),
		log.String("app", "start ok"),
	)
}

func customStdLogger() {
	log.ResetDefault(log.New(os.Stdout, log.DebugLevel, log.WithCaller(false)))
	defer log.Sync()

	log.Debug("custom stdout logger",
		log.String("level", "debug"),
		log.String("app", "start ok"),
	)
}

func main() {
	defaultLogger()
	fmt.Println("--------------")
	customStdLogger()
}
