package file

import (
	"io"

	"github.com/jfindley/automate/core"
)

// File implements a Module which operates on files
type File struct {
	path    string
	content io.Reader
}

func (f *File) Name() string {
	return "file"
}

func (f *File) Configure(in core.Input) error {
	err := in.Validate(schema)
	if err != nil {
		return err
	}

	f.path = in.Data()["path"].(string)

	return err
}

func (f *File) Run(r core.ResponseWriter) {

}
