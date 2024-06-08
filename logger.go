package grimoire

import "fmt"

type logger interface {
	Error(err error)
	Debug(string)
}

func newdefaultLogger() logger {
	return &defaultLogger{}
}

type defaultLogger struct{}

func (d defaultLogger) Error(err error) {
	fmt.Println(err)
}

func (d defaultLogger) Debug(s string) {
	fmt.Println(s)
}
