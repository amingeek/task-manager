// backend/utils/logger.go

package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func getLogLevelString(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func logMessage(level LogLevel, message string, details ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := getLogLevelString(level)
	
	logEntry := fmt.Sprintf("[%s] %s: %s", timestamp, levelStr, message)
	
	if len(details) > 0 {
		logEntry += fmt.Sprintf(" | %v", details)
	}
	
	// نوشتن در فایل و console
	log.Println(logEntry)
	
	// نوشتن در فایل لاگ
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		defer file.Close()
		fmt.Fprintln(file, logEntry)
	}
}

func LogDebug(message string, details ...interface{}) {
	logMessage(LogLevelDebug, message, details...)
}

func LogInfo(message string, details ...interface{}) {
	logMessage(LogLevelInfo, message, details...)
}

func LogWarn(message string, details ...interface{}) {
	logMessage(LogLevelWarn, message, details...)
}

func LogError(message string, err error) {
	logMessage(LogLevelError, message, err.Error())
}

func LogErrorWithDetails(message string, err error, details ...interface{}) {
	logMessage(LogLevelError, message, fmt.Sprintf("%v | Details: %v", err, details))
}
