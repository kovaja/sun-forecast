package logger

import (
	"fmt"
	"time"
)

func getAppPrefix() string {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z07:00")
	return fmt.Sprintf("[Sun-forecast][%s]", timestamp)
}

func Log(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s: %s\n", getAppPrefix(), msg)
}

func LogError(msg string, err error) {
	fmt.Printf(
		"%s: %s. Error: %v\n",
		getAppPrefix(),
		msg,
		err,
	)
}
