package file

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jfindley/automate/core"
	"github.com/stretchr/testify/assert"
)

var testFile = os.TempDir() + "/automate_file"
var testData = []byte("test data")
var testInitialData = []byte("initial data")

func TestFileName(t *testing.T) {
	f := new(File)

	if f.Name() != "file" {
		t.Error("Wrong module name")
	}
}

func TestFileConfigure(t *testing.T) {
	f := new(File)

	conf := core.NewConfigInput(map[string]interface{}{})

	err := f.Configure(conf)
	if err == nil {
		t.Error("No error with missing parameters")
	}

	conf = core.NewConfigInput(map[string]interface{}{
		"path": testFile,
	})
	err = f.Configure(conf)
	assert.Nil(t,err)
	if f.path != testFile {
		t.Error("Wrong file path")
	}
	if f.mode != os.FileMode(0644) {
		t.Error("Wrong file mode")
	}

	conf = core.NewConfigInput(map[string]interface{}{
		"path": testFile,
		"mode": "bad",
	})
	err = f.Configure(conf)
	if err.Error() != "Unable to parse mode" {
		t.Error("Bad error status")
	}

	conf = core.NewConfigInput(map[string]interface{}{
		"path": testFile,
		"mode": "0755",
	})
	err = f.Configure(conf)
	assert.Nil(t,err)

	if f.mode != os.FileMode(0755) {
		t.Error("Wrong file mode", f.mode)
	}

	conf = core.NewConfigInput(map[string]interface{}{
		"path":   testFile,
		"mode":   0755,
		"action": "touch",
	})
	err = f.Configure(conf)
	assert.Nil(t,err)

	if f.mode != os.FileMode(0755) {
		t.Error("Wrong file mode", f.mode)
	}

}

func TestFileTouch(t *testing.T) {
	f := new(File)

	conf := core.NewConfigInput(map[string]interface{}{
		"path":   testFile,
		"mode":   0755,
		"action": "touch",
	})
	err := f.Configure(conf)
	assert.Nil(t,err)

	err = f.action()
	assert.Nil(t,err)
	defer os.Remove(testFile)

	fi, err := os.Stat(testFile)
	assert.Nil(t,err)

	if fi.Mode() != os.FileMode(0755) {
		t.Error("Wrong file mode", f.mode)
	}
}

func TestFileRemove(t *testing.T) {
	f := new(File)

	file, err := os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY, f.mode)
	assert.Nil(t,err)
	file.Close()

	conf := core.NewConfigInput(map[string]interface{}{
		"path":   testFile,
		"action": "remove",
	})
	err = f.Configure(conf)
	assert.Nil(t,err)

	err = f.action()
	assert.Nil(t,err)

	_, err = os.Stat(testFile)
	if !os.IsNotExist(err) {
		t.Error("File not removed")
	}
}

func TestSetFile(t *testing.T) {
	f := new(File)

	err := ioutil.WriteFile(testFile, testInitialData, 0644)
	assert.Nil(t,err)
	defer os.Remove(testFile)

	conf := core.NewConfigInput(map[string]interface{}{
		"path":   testFile,
		"action": "set",
	})
	err = f.Configure(conf)
	assert.Nil(t,err)

	f.data = bytes.NewReader(testData)
	err = f.action()
	assert.Nil(t,err)

	res, err := ioutil.ReadFile(testFile)
    assert.Nil(t, err)
    
    assert.Equal(t, testData, res, "File contents should match")

}
