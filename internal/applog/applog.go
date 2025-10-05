// Description: This package provides log management for the application.
// Developer: Aleksei Grigorev <https://github.com/AlekseiGrigorev>, <aleksvgrig@gmail.com>
// Copyright (c) 2025 Aleksei Grigorev
package applog

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Constants for date and time formatting.
const (
	REFERENCE_DATE = "2006-01-02 15:04:05"
)

// AppLog represents a logging structure with file output capabilities and error tracking.
// It allows logging to a file and tracks whether any write errors have occurred.
type AppLog struct {
	// File to which log messages will be written.
	// If nil, no file output will occur.
	File *os.File
	// Id is a unique identifier for the log entry.
	Id string
	// Mutex is used to synchronize color output.
	// color is non-thread safe, so we need to synchronize it.
	Mutex *sync.Mutex
	// hasWriteFileError tracks whether any write errors have occurred during file writing.
	// If true, subsequent attempts to write to the file will be skipped.
	hasWriteFileError bool
}

// getDate returns the current date and time as a string in the format "YYYY-MM-DD HH:MM:SS".
func (appLog *AppLog) getDate() string {
	return time.Now().Format(REFERENCE_DATE)
}

// insertDate prepends the current date and time, followed by a hyphen,
// to the provided arguments. It returns the modified argument list
// which includes the date and time as the first element.
func (appLog *AppLog) insertDateAndId(args ...any) []any {
	if appLog.Id != "" {
		args = append([]any{appLog.getDate(), "-", appLog.Id, "-"}, args...)
	} else {
		args = append([]any{appLog.getDate(), "-"}, args...)
	}
	return args
}

// writeToFile writes the given string to the file specified by the File field.
// If the file is nil or hasWriteFileError is true, the method does nothing.
// If the write operation fails, the method sets hasWriteFileError to true and
// logs the error.
func (appLog *AppLog) writeToFile(str string) (int int, err error) {
	if appLog.File == nil || appLog.hasWriteFileError {
		return
	}
	res, err := appLog.File.WriteString(str)
	if err != nil {
		appLog.hasWriteFileError = true
		appLog.Error("Error writing to log file:", err)
	}
	return res, err
}

// String formats the provided arguments by prepending the current date and time,
// followed by a hyphen, and returns the formatted string.
func (appLog *AppLog) String(args ...any) string {
	return fmt.Sprintln(appLog.insertDateAndId(args...)...)
}

// Info logs the given arguments as informational messages by prepending "[Info]" to the
// arguments, formatting them with the current date and time, and writing the formatted
// string to the file specified by the File field. The method returns itself.
func (appLog *AppLog) Info(args ...any) *AppLog {
	args = append([]any{"[Info]"}, args...)
	str := appLog.String(args...)
	appLog.coloredLog(str, color.New(color.FgBlue))
	appLog.writeToFile(str)
	return appLog
}

// Ok logs the given arguments as success messages by prepending "[Ok]" to the
// arguments, formatting them with the current date and time, and writing the
// formatted string to the file specified by the File field. The method returns
// itself.
func (appLog *AppLog) Ok(args ...any) *AppLog {
	args = append([]any{"[Ok]"}, args...)
	str := appLog.String(args...)
	appLog.coloredLog(str, color.New(color.FgGreen))
	appLog.writeToFile(str)
	return appLog
}

// Warn logs the given arguments as warning messages by prepending "[Warn]" to the
// arguments, formatting them with the current date and time, and writing the
// formatted string to the file specified by the File field. The method returns
// itself.
func (appLog *AppLog) Warn(args ...any) *AppLog {
	args = append([]any{"[Warn]"}, args...)
	str := appLog.String(args...)
	appLog.coloredLog(str, color.New(color.FgYellow))
	appLog.writeToFile(str)
	return appLog
}

// Error logs the given arguments as error messages by prepending "[Error]" to the
// arguments, formatting them with the current date and time, and writing the
// formatted string to the file specified by the File field. The method returns
// itself.
func (appLog *AppLog) Error(args ...any) *AppLog {
	args = append([]any{"[Error]"}, args...)
	str := appLog.String(args...)
	appLog.coloredLog(str, color.New(color.FgRed))
	appLog.writeToFile(str)
	return appLog
}

// Write logs the given arguments as a message of the given type by prepending the appropriate
// prefix to the arguments, formatting them with the current date and time, and writing the
// formatted string to the file specified by the File field. The method returns itself.
// The type of the message can be one of the following: "info", "ok", "warn", "error".
// If the type is not recognized, the method defaults to logging the message as an informational message.
func (appLog *AppLog) Write(msgType string, args ...any) *AppLog {
	switch msgType {
	case "info":
		return appLog.Info(args...)
	case "ok":
		return appLog.Ok(args...)
	case "warn":
		return appLog.Warn(args...)
	case "error":
		return appLog.Error(args...)
	default:
		return appLog.Info(args...)
	}
}

// coloredLog prints the given string with the given color to the console.
// If a Mutex is set for the AppLog, it locks the mutex before printing the string
// to ensure thread safety. It also sets color.NoColor to false before printing so
// that the color is actually applied. If no Mutex is set, it simply prints the
// string with the given color.
func (appLog *AppLog) coloredLog(str string, newColor *color.Color) *AppLog {
	if appLog.Mutex != nil {
		appLog.Mutex.Lock()
		color.NoColor = false
		newColor.Print(str)
		appLog.Mutex.Unlock()
	} else {
		newColor.Println(str)
	}
	return appLog
}
