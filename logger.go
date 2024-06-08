package grimoire

import (
	"log"
)

type logger interface {
	Error(err error)
	Debug(string)
}

func newdefaultLogger() logger {
	return &defaultLogger{}
}

type defaultLogger struct{}

func (d defaultLogger) Error(err error) {
	log.Fatalln(err)
}

func (d defaultLogger) Debug(s string) {
	log.Println(s)
}
