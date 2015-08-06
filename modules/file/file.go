package file

import (
	"io"
	"os"

	"github.com/jfindley/automate/core"
)

// File implements a Module which operates on files
type File struct {
	handle *os.File
	ctl    core.RunControl
	source io.Reader
}

func (f *File) Name() string {
	return "file"
}

func (f *File) Configure(in core.Input) error {
	err := in.Validate(schema)

	return err
}

func (f *File) Run(r core.ResponseWriter) {

}
