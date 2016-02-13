package file

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/jfindley/automate/core"
	"github.com/jfindley/testfs"
)

var fs testfs.FileSystem

func init() {
	fs = testfs.NewOSFS()
}

// File implements a Module which operates on files.
type File struct {
	path   string
	mode   os.FileMode
	data   io.Reader
	action func() error
}

func (f *File) Name() string {
	return "file"
}

func (f *File) Configure(in core.Input) error {
	err := in.Validate(schema)
	if err != nil {
		return err
	}
    
    var ok bool

	f.path, ok = in.Data()["path"].(string)
    if !ok {
        return errors.New("Path is not a string")
    }

	if val, ok := in.Data()["mode"]; !ok {
		f.mode = os.FileMode(0644)
	} else {
		f.mode, err = fileMode(val)
		if err != nil {
			return errors.New("Unable to parse mode")
		}
	}

	if val, ok := in.Data()["action"]; !ok {
		f.action = f.touch
	} else {
		switch val.(string) {
		case "touch":
			f.action = f.touch

		case "remove":
			f.action = f.remove

		case "set":
			f.action = f.set

		default:
			return errors.New("Unable to parse action")
		}
	}

	return err
}

func (f *File) Run(r core.ResponseWriter) {
	// var origExist bool
	// fi, err := os.Stat(f.path)
}

func (f *File) RunTest(testFs testfs.FileSystem, r core.ResponseWriter) {
	fs = testFs
	f.Run(r)
}

// touch a file.  Same as system touch.
func (f *File) touch() error {
	file, err := fs.OpenFile(f.path, os.O_CREATE|os.O_WRONLY, f.mode)
	if err != nil {
		return err
	}
	return file.Close()
}

// set sets the contents of a file.
func (f *File) set() error {
	file, err := fs.OpenFile(f.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.mode)
	if err != nil {
		return err
	}
	defer file.Close()

	return bufferedWrite(f.data, file)
}

// remove removes a file.
func (f *File) remove() error {
	return fs.Remove(f.path)
}

// fileMode parses input into a valid file mode
func fileMode(in interface{}) (os.FileMode, error) {

	switch in.(type) {
	case os.FileMode:
		return in.(os.FileMode), nil

	// We have special handling for strings, as the default base-10 handling
	// in ConvertInt will produce incorrect results
	case string:
		m, err := strconv.ParseInt(in.(string), 8, 32)
		return os.FileMode(m), err

	default:
		m, err := core.ConvertInt(in)
		return os.FileMode(m), err
	}

}

// bufferedWrite efficiently writes from an io.Reader to an io.Writer
func bufferedWrite(in io.Reader, out io.Writer) error {
	w := bufio.NewWriter(out)

	_, err := io.Copy(w, in)
	if err != nil {
		return err
	}

	return w.Flush()
}
