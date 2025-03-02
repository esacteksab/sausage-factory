package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/log"
)

var (
	butter bool
)

func main() {
	log.Debug("Cookie :cookie:")
	log.Info("Hello World")

	err := fmt.Errorf("too much sugar")
	log.Error("Failed to bake cookies", "err", err)

	log.Print("Baking 101")
	// 2023/01/04 10:04:06 Baking 101

	logger := log.New(os.Stderr)
	if butter {
		logger.Warn("chewy!", "butter", true)
	}

	// format := "%s %d"
	// log.Debugf(format, "chocolate", 10)
	// log.Warnf(format, "adding more", 5)
	// log.Errorf(format, "increasing temp", 420)
	// log.Fatalf(format, "too hot!", 500)  // this calls os.Exit(1)
	// log.Printf(format, "baking cookies") // prints regardless of log level

	// Use these in conjunction with `With(...)` to add more context
	log.With("err", err).Errorf("unable to start %s", "oven")

	//newlogger := log.New(os.Stderr)
	//newlogger.SetReportTimestamp(false)
	//newlogger.SetReportCaller(false)
	//newlogger.SetLevel(log.DebugLevel)

	handler := log.New(os.Stderr)
	slogger := slog.New(handler)
	slogger.Error("meow?")
}
