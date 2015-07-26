package core

import (
	"github.com/Sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
}

type Response struct {
	Ok        bool
	Callbacks bool
	Fields    logrus.Fields
}

func (r *Response) Success(ok bool) {
	r.Ok = ok
}

func (r *Response) TriggerCallbacks(ok bool) {
	r.Callbacks = ok
}

func (r *Response) Message(level string, messages ...interface{}) {
	l := log.WithFields(r.Fields)

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		l.Warn("Unable to parse log level: ", level)
		lvl = logrus.InfoLevel
	}

	switch lvl {
	case logrus.PanicLevel:
		l.Panic(messages...)
	case logrus.FatalLevel:
		l.Fatal(messages...)
	case logrus.ErrorLevel:
		l.Error(messages...)
	case logrus.WarnLevel:
		l.Warn(messages...)
	case logrus.InfoLevel:
		l.Info(messages...)
	case logrus.DebugLevel:
		l.Debug(messages...)
	}
}

func NewResponse() *Response {
	return new(Response)
}
