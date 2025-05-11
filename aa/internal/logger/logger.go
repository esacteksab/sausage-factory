// SPDX-License-Identifier: MIT
package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	Logger  *log.Logger
	Verbose bool
)

func init() {
	SetVerbose(Verbose)
}

// initLogger is how we initially create Logger. The values passed are based on
// 'Verbose' being true Colors are defined here
// https://github.com/charmbracelet/x/blob/aedd0cd23ed703ff7cbccc5c4f9ab51a4768a9e6/ansi/color.go#L15-L32
// 14 is Bright Cyan, 9 is Red -- no more purple
func setupLogger(ReportCaller, ReportTimestamp bool, TimeFormat string) (Logger *log.Logger) {
	Logger = log.NewWithOptions(
		os.Stderr, log.Options{
			ReportCaller:    ReportCaller,
			ReportTimestamp: ReportTimestamp,
			TimeFormat:      TimeFormat,
		},
	)
	MaxWidth := 4
	styles := log.DefaultStyles()
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString(strings.ToUpper(log.DebugLevel.String())).
		Bold(true).MaxWidth(MaxWidth).Foreground(lipgloss.Color("14"))

	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString(strings.ToUpper(log.FatalLevel.String())).
		Bold(true).MaxWidth(MaxWidth).Foreground(lipgloss.Color("9"))
	Logger.SetStyles(styles)
	Logger.SetLevel(log.DebugLevel)
	log.SetDefault(Logger)
	return Logger
}

func SetVerbose(verbose bool) {
	Verbose = verbose

	log.Debugf("Verbose is %v\n", Verbose)
	log.Debugf("verbose is %v\n", verbose)
	if Verbose {
		Logger = setupLogger(true, true, "2006/01/02 15:04:05")
		log.SetLevel(log.DebugLevel)
		log.SetDefault(Logger)
	} else {
		Logger = setupLogger(false, false, time.Kitchen)
		log.SetLevel(log.InfoLevel)
		log.SetDefault(Logger)
	}
}

// Info logs information
func Info(msg string, keyvals ...interface{}) {
	Logger.Info(msg, keyvals...)
}

// Infof logs formatted information
func Infof(format string, args ...interface{}) {
	Logger.Info(fmt.Sprintf(format, args...))
}

// Error logs errors
func Error(msg string, keyvals ...interface{}) {
	Logger.Error(msg, keyvals...)
}

// Errorf logs formatted errors
func Errorf(format string, args ...interface{}) {
	Logger.Error(fmt.Sprintf(format, args...))
}

// Fatal logs fatal errors and exits
func Fatal(msg string, keyvals ...interface{}) {
	Logger.Fatal(msg, keyvals...)
}

// Fatalf logs formatted fatal errors and exits
func Fatalf(format string, args ...interface{}) {
	Logger.Fatal(fmt.Sprintf(format, args...))
}

// Debug logs debug information
func Debug(msg string, keyvals ...interface{}) {
	Logger.Debug(msg, keyvals...)
}

// Debugf logs formatted debug information
func Debugf(format string, args ...interface{}) {
	Logger.Debug(fmt.Sprintf(format, args...))
}

// Warn logs warnings
func Warn(msg string, keyvals ...interface{}) {
	Logger.Warn(msg, keyvals...)
}

// Warnf logs formatted warnings
func Warnf(format string, args ...interface{}) {
	Logger.Warn(fmt.Sprintf(format, args...))
}

// SetLevel sets the logging level
func SetLevel(level log.Level) {
	Logger.SetLevel(level)
}
