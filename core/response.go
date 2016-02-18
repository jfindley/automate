package core

import (
	log "github.com/Sirupsen/logrus"
)

type Response struct {
	Ok     bool
	Notify bool
	Fields log.Fields
}

func (r *Response) Success(ok bool) {
	r.Ok = ok
}

func (r *Response) Changed(ok bool) {
	r.Notify = ok
}

func (r *Response) Message(level string, messages ...interface{}) {
	l := log.WithFields(r.Fields)

	lvl, err := log.ParseLevel(level)
	if err != nil {
		l.Warn("Unable to parse log level: ", level)
		lvl = log.InfoLevel
	}

	switch lvl {
	case log.PanicLevel:
		l.Panic(messages...)
	case log.FatalLevel:
		l.Fatal(messages...)
	case log.ErrorLevel:
		l.Error(messages...)
	case log.WarnLevel:
		l.Warn(messages...)
	case log.InfoLevel:
		l.Info(messages...)
	case log.DebugLevel:
		l.Debug(messages...)
	}
}

func NewResponse() *Response {
	return new(Response)
}
