package file

import (
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
	assert.Error(t, err, "No error with missing parameters")

	conf = core.NewConfigInput(map[string]interface{}{
		"path": testFile,
	})
	err = f.Configure(conf)
	assert.NoError(t, err)

	assert.Equal(t, testFile, f.path, "File path should match")
	assert.Equal(t, os.FileMode(0644), f.mode, "File mode should match")

	conf = core.NewConfigInput(map[string]interface{}{
		"path": testFile,
		"mode": "bad",
	})
	err = f.Configure(conf)
	assert.EqualError(t, err, "Unable to parse mode")

	conf = core.NewConfigInput(map[string]interface{}{
		"path": testFile,
		"mode": "0755",
	})
	err = f.Configure(conf)
	assert.NoError(t, err)

	assert.Equal(t, os.FileMode(0755), f.mode, "File mode should match")

	conf = core.NewConfigInput(map[string]interface{}{
		"path":   testFile,
		"mode":   0755,
		"action": "touch",
	})
	err = f.Configure(conf)
	assert.NoError(t, err)

	assert.Equal(t, os.FileMode(0755), f.mode, "File mode should match")

}

func TestFileTouch(t *testing.T) {
	f := new(File)

	conf := core.NewConfigInput(map[string]interface{}{
		"path":   testFile,
		"mode":   0755,
		"action": "touch",
	})
	err := f.Configure(conf)
	assert.NoError(t, err)

	err = f.action()
	assert.NoError(t, err)
	defer os.Remove(testFile)

	fi, err := os.Stat(testFile)
	assert.NoError(t, err)

	assert.Equal(t, os.FileMode(0755), fi.Mode(), "File mode should match")
}

func TestFileRemove(t *testing.T) {
	f := new(File)

	file, err := os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY, f.mode)
	assert.NoError(t, err)
	file.Close()

	conf := core.NewConfigInput(map[string]interface{}{
		"path":   testFile,
		"action": "remove",
	})
	err = f.Configure(conf)
	assert.NoError(t, err)

	err = f.action()
	assert.NoError(t, err)

	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))
}

func TestSetFile(t *testing.T) {
	f := new(File)

	err := ioutil.WriteFile(testFile, testInitialData, 0644)
	assert.NoError(t, err)
	defer os.Remove(testFile)

	conf := core.NewConfigInput(map[string]interface{}{
		"path":    testFile,
		"action":  "set",
		"content": testData,
	})
	err = f.Configure(conf)
	assert.NoError(t, err)

	err = f.action()
	assert.NoError(t, err)

	res, err := ioutil.ReadFile(testFile)
	assert.NoError(t, err)

	assert.Equal(t, testData, res, "File contents should match")

}
