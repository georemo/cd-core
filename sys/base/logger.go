package base

import (
	"fmt"

	"github.com/fatih/color"
)

// Logger represents a logger with different log levels
type Logger struct{}

// LogWarning logs a warning message with yellow color
func (l *Logger) LogWarning(message string) {
	warn := color.New(color.FgYellow).SprintFunc()
	fmt.Println("[WARNING]", warn(message))
}

// LogError logs an error message with red color
func (l *Logger) LogError(message string) {
	error := color.New(color.FgRed).SprintFunc()
	fmt.Println("[ERROR]", error(message))
}

// LogInfo logs an info message with green color
func (l *Logger) LogInfo(message string) {
	info := color.New(color.FgGreen).SprintFunc()
	fmt.Println("[INFO]", info(message))
}
