package logger

import (
	"fmt"
	"github.com/cgCodeLife/image_recognition_web/config"
	"github.com/cgCodeLife/logs"
	"os"
	"path/filepath"
)

const (
	MAX_LOG_SIZE = 1024 * 1024 * 1024
)

var (
	webLogger *logs.Logger
)

func Init() {
	initAppLogger()
	logs.InitLogger(webLogger)
}

func initAppLogger() {
	level := config.Level()
	webLogger = logs.NewLogger(1024)
	webLogger.SetLevel(level)
	webLogger.SetPSM(config.PSM())
	webLogger.SetCallDepth(3)
	if config.FileLog() {
		webLog := filepath.Join(config.LogDir(), config.PSM()+".log")
		fileProvider := logs.NewFileProvider(webLog, logSegment(config.LogInterval()), MAX_LOG_SIZE)
		fileProvider.SetLevel(level)
		if err := webLogger.AddProvider(fileProvider); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to add fileProvider: %s\n", err)
		}
	}
	if config.ConsoleLog() {
		consoleProvider := logs.NewConsoleProvider()
		consoleProvider.SetLevel(level)
		if err := webLogger.AddProvider(consoleProvider); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to add consoleProvider: %s\n", err)
		}
	}
}

func Stop() {
	if webLogger != nil {
		webLogger.Stop()
	}
}

func logSegment(logInterval string) logs.SegDuration {
	switch logInterval {
	case "hour":
		return logs.HourDur
	case "day":
		return logs.DayDur
	default:
	}
	return logs.NoDur
}
