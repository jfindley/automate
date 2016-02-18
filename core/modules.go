package core

import (
	"sourcegraph.com/sourcegraph/rwvfs"
)

type Module interface {
	Name() string
	Configure(Input) error
	Run(rwvfs.FileSystem, ResponseWriter)
}

type ResponseWriter interface {
	Success(bool)
	Changed(bool)
	Message(string, ...interface{})
}

type Input interface {
	Data(string) (interface{}, error)
	Type(string) (string, error)
	Validate(InputSchema) error
}
