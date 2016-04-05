package file

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"io/ioutil"

	"github.com/jfindley/automate/core"
	"github.com/jfindley/testfs"
)

// File implements a Module which operates on files.
type File struct {
	path   string
	mode   os.FileMode
	data   io.Reader
	sum    []byte
	action func(testfs.FileSystem, core.ResponseWriter)
}

// Name returns the name of the module
func (f *File) Name() string {
	return "file"
}

// Configure configures an instance of a module
func (f *File) Configure(in core.Input) error {
	err := in.Validate(schema)
	if err != nil {
		return err
	}

	if val, err := in.Data("path"); err != nil {
		return errors.New("Path is required")
	} else {
		f.path = val.(string)
	}

	if val, err := in.Data("mode"); err != nil {
		f.mode = os.FileMode(0644)
	} else {
		f.mode, err = fileMode(val)
		if err != nil {
			return errors.New("Unable to parse mode")
		}
	}

	if val, err := in.Data("action"); err != nil {
		f.action = f.touch
	} else {
		switch val.(string) {
		case "touch":
			f.action = f.touch

		case "remove":
			f.action = f.remove

		case "set":
			f.configureData(in)
			f.action = f.set

		default:
			return errors.New("Unable to parse action")
		}
	}

	return err
}

func (f *File) configureData(in core.Input) {

	t, err := in.Type("content")
	if err != nil {
		return
	}

	if t == "pipe" {
		val, _ := in.Data("content")
		f.data = val.(io.Reader)
	}

	if t == "data" {
		val, _ := in.Data("content")
		f.sum = sum(val.([]byte))
		f.data = bytes.NewReader(val.([]byte))
	}
}

// Run executes a module instance
func (f *File) Run(fs testfs.FileSystem, r core.ResponseWriter) {
	f.action(fs, r)
}

// touch a file.  Same as system touch.
func (f *File) touch(fs testfs.FileSystem, r core.ResponseWriter) {
	file, err := fs.OpenFile(f.path, os.O_CREATE|os.O_WRONLY, f.mode)
	if err != nil {
		r.Message("error", err.Error())
		r.Success(false)
		return
	}
	err = file.Close()
	if err != nil {
		r.Message("error", err.Error())
		r.Success(false)
		return
	}
	r.Message("info", "touched ", f.path)
	r.Success(true)
	r.Changed(true)
	return
}

// set sets the contents of a file.
func (f *File) set(fs testfs.FileSystem, r core.ResponseWriter) {
	var (
		file testfs.File
		err  error
	)

	// If we have a valid checksum for the input data and the file exists, read it
	// and avoid modifying the file if possible.
	if f.matching(fs) {
		r.Success(true)
		r.Changed(false)
		r.Message("info", f.path, " unchanged")
		return
	}

	file, err = fs.OpenFile(f.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.mode)
	if err != nil {
		r.Message("error", err.Error())
		r.Success(false)
		return
	}

	err = bufferedWrite(f.data, file)
	if err != nil {
		r.Message("error", err.Error())
		r.Success(false)
		return
	}

	err = file.Close()
	if err != nil {
		r.Message("error", err.Error())
		r.Success(false)
		return
	}

	r.Message("info", "set content of ", f.path)
	r.Success(true)
	r.Changed(true)
	return
}

// remove removes a file.
func (f *File) remove(fs testfs.FileSystem, r core.ResponseWriter) {
	err := fs.Remove(f.path)
	switch {

	// Don't error if the file is already removed
	case os.IsNotExist(err):
		r.Success(true)
		r.Changed(false)
		r.Message("info", f.path, " already removed")

	case err == nil:
		r.Success(true)
		r.Changed(true)
		r.Message("info", f.path, " removed")

	default:
		r.Message("error", err.Error())
		r.Success(false)

	}

	return
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
	if in == nil {
		return errors.New("Input cannot be nil")
	}
	if out == nil {
		return errors.New("Output cannot be nil")
	}

	w := bufio.NewWriter(out)

	_, err := io.Copy(w, in)
	if err != nil {
		return err
	}

	return w.Flush()
}

// matching returns true if a file matches the checksum, false otherwise
func (f *File) matching(fs testfs.FileSystem) bool {
	if f.sum == nil || len(f.sum) == 0 {
		return false
	}
	_, err := fs.Stat(f.path)
	if err != nil {
		return false
	}
	file, err := fs.OpenFile(f.path, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}

	if bytes.Compare(f.sum, data) == 0 {
		return true
	}

	return false
}
