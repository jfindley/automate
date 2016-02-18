package file

import (
	"os"
	"testing"
    "bytes"

	"github.com/jfindley/automate/core"
	"github.com/jfindley/testfs"
	"github.com/stretchr/testify/assert"
	log "github.com/Sirupsen/logrus"
)

var resp *core.Response

func init() {
    resp = core.NewResponse()
    var out bytes.Buffer
    log.SetOutput(&out)
}

var testFile = "/test_file"
var testData = []byte("test data")
var testInitialData = []byte("initial data")

func TestFileName(t *testing.T) {
	file := new(File)
    
    assert.Equal(t, "file", file.Name())
}

func TestFileConfigure(t *testing.T) {
	file := new(File)

	conf := core.NewConfigInput(map[string]interface{}{})

	err := file.Configure(conf)
	assert.Error(t, err, "No error with missing parameters")

	conf = core.NewConfigInput(map[string]interface{}{
		"path": testFile,
	})
	err = file.Configure(conf)
	assert.NoError(t, err)

	assert.Equal(t, testFile, file.path, "File path should match")
	assert.Equal(t, os.FileMode(0644), file.mode, "File mode should match")

	conf = core.NewConfigInput(map[string]interface{}{
		"path": testFile,
		"mode": "bad",
	})
	err = file.Configure(conf)
	assert.EqualError(t, err, "Unable to parse mode")

	conf = core.NewConfigInput(map[string]interface{}{
		"path": testFile,
		"mode": "0755",
	})
	err = file.Configure(conf)
	assert.NoError(t, err)

	assert.Equal(t, os.FileMode(0755), file.mode, "File mode should match")

	conf = core.NewConfigInput(map[string]interface{}{
		"path":   testFile,
		"mode":   0755,
		"action": "touch",
	})
	err = file.Configure(conf)
	assert.NoError(t, err)

	assert.Equal(t, os.FileMode(0755), file.mode, "File mode should match")

}

func TestFileTouch(t *testing.T) {
	file := new(File)
	fs := testfs.NewLocalTestFS()

	conf := core.NewConfigInput(map[string]interface{}{
		"path":   testFile,
		"mode":   0755,
		"action": "touch",
	})
	err := file.Configure(conf)
	assert.NoError(t, err)

	resp := core.NewResponse()

	file.Run(fs, resp)
	assert.True(t, resp.Ok)
	defer os.Remove(testFile)

	fi, err := fs.Stat(testFile)
	assert.NoError(t, err)

	assert.Equal(t, os.FileMode(0755), fi.Mode(), "File mode should match")
}

func TestFileRemove(t *testing.T) {
	file := new(File)
	fs := testfs.NewLocalTestFS()

	_, err := fs.Create(testFile)
	assert.NoError(t, err)

	conf := core.NewConfigInput(map[string]interface{}{
		"path":   testFile,
		"action": "remove",
	})
	err = file.Configure(conf)
	assert.NoError(t, err)

	resp := core.NewResponse()

	file.Run(fs, resp)
	assert.True(t, resp.Ok)

	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))
}

func TestSetFile(t *testing.T) {
	file := new(File)
	fs := testfs.NewLocalTestFS()

	_, err := fs.Create(testFile)
	assert.NoError(t, err)

	conf := core.NewConfigInput(map[string]interface{}{
		"path":    testFile,
		"action":  "set",
		"content": testData,
	})
	err = file.Configure(conf)
	assert.NoError(t, err)

	resp := core.NewResponse()

	file.Run(fs, resp)
	assert.True(t, resp.Ok)

	res, err := fs.Open(testFile)
	assert.NoError(t, err)

	out := make([]byte, len(testData))
	_, err = res.Read(out)
    assert.NoError(t, err)

	assert.Equal(t, testData, out, "File contents should match")

}
