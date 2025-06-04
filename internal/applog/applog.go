package applog

import (
	"fmt"

	"github.com/fatih/color"
)

type AppLog struct {
}

func (appLog *AppLog) Info(args ...any) *AppLog {
	fmt.Println(args...)
	return appLog
}

func (appLog *AppLog) Ok(args ...any) *AppLog {
	color.NoColor = false
	color.Green(fmt.Sprintln(args...))
	return appLog
}

func (appLog *AppLog) Warn(args ...any) *AppLog {
	color.NoColor = false
	color.Yellow(fmt.Sprintln(args...))
	return appLog
}

func (appLog *AppLog) Error(args ...any) *AppLog {
	color.NoColor = false
	color.Red(fmt.Sprintln(args...))
	return appLog
}
