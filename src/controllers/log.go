package controllers

import (
	"fmt"
	"log"
)

type Logger struct {
	StartWith string
}

func NewLogger() *Logger {
	return &Logger{}
}
func (l *Logger) Error(err error, msg string, KeyAndValues ...interface{}) {
	ErrMsg := l.StartWith + fmt.Sprintf("Error: %v, Message: %v.\n", err, msg)
	for i := 0; KeyAndValues != nil && i < len(KeyAndValues); i += 2 {
		ErrMsg += fmt.Sprintf("	%v:%v.\n", KeyAndValues[i], KeyAndValues[i+1])
	}
	log.Println(ErrMsg)
}

func (l *Logger) Info(msg string, KeyAndValues ...interface{}) {
	ErrMsg := l.StartWith + fmt.Sprintf("Info: %v.\n", msg)

	for i := 0; KeyAndValues != nil && i < len(KeyAndValues); i += 2 {
		ErrMsg += fmt.Sprintf("	%v:%v.\n", KeyAndValues[i], KeyAndValues[i+1])
	}
	log.Println(ErrMsg)
}

func (l *Logger) WithName(name string) *Logger {
	l.StartWith += name
	l.StartWith += ": "
	return l
}
