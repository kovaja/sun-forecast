package logger

import "fmt"

func getAppPrefix() string {
	return "[Sun-forecast]"
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
