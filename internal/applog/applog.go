package applog

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

const (
	REFERENCE_DATE = "2006-01-02 15:04:05"
)

type AppLog struct {
	File              *os.File
	hasWriteFileError bool
}

func (appLog *AppLog) getDate() string {
	return time.Now().Format(REFERENCE_DATE)
}

func (appLog *AppLog) insertDate(args ...any) []any {
	args = append([]any{appLog.getDate(), "-"}, args...)
	return args
}

func (appLog *AppLog) writeToFile(str string) (int int, err error) {
	if appLog.File == nil || appLog.hasWriteFileError {
		return
	}
	res, err := appLog.File.WriteString(str)
	if err != nil {
		appLog.Error("Error writing to log file:", err)
		appLog.hasWriteFileError = true
	}
	return res, err
}

func (appLog *AppLog) String(args ...any) string {
	return fmt.Sprintln(appLog.insertDate(args...)...)
}

func (appLog *AppLog) Info(args ...any) *AppLog {
	args = append([]any{"[Info]"}, args...)
	str := appLog.String(args...)
	color.NoColor = false
	color.Blue(str)
	appLog.writeToFile(str)
	return appLog
}

func (appLog *AppLog) Ok(args ...any) *AppLog {
	args = append([]any{"[Ok]"}, args...)
	str := appLog.String(args...)
	color.NoColor = false
	color.Green(str)
	appLog.writeToFile(str)
	return appLog
}

func (appLog *AppLog) Warn(args ...any) *AppLog {
	args = append([]any{"[Warn]"}, args...)
	str := appLog.String(args...)
	color.NoColor = false
	color.Yellow(str)
	appLog.writeToFile(str)
	return appLog
}

func (appLog *AppLog) Error(args ...any) *AppLog {
	args = append([]any{"[Error]"}, args...)
	str := appLog.String(args...)
	color.NoColor = false
	color.Red(str)
	appLog.writeToFile(str)
	return appLog
}
