package main

import (
	"fmt"
	"os"
)

type Logger struct {
	ClientName string
	FilePath   string
}

func NewLogger(ClientName string, FilePath string) *Logger {
	var logger Logger
	logger.ClientName = ClientName
	logger.FilePath = FilePath

	return &logger
}

// creates a new log and adds to file. Generates an alert email if alert is true
func (l *Logger) Log(alert_type string, message string, alert bool) error {
	//adds information to the log message
	formattedLog := fmt.Sprintf("[%s] %s: %s", l.ClientName, alert_type, message)

	//adds this log to the file
	file, err := os.OpenFile(l.FilePath, os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("unable to open logging file: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(formattedLog)
	if err != nil {
		return fmt.Errorf("unable to write to logging file %s", err)
	}

	//send alert if needed
	if alert {
		l.Alert(formattedLog)
	}
	return nil
}

func (l *Logger) Alert(message string) error {
	//sends an email alert based on this message

	return nil
}
