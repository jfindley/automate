package docker

import (
	"testing"

	"github.com/jfindley/testfs"
	"github.com/stretchr/testify/assert"
)

func TestNewLayer(t *testing.T) {

	fs := testfs.NewLocalTestFS()

	err := fs.Mkdir("/test", 0755)
	assert.NoError(t, err)

	file, err := fs.Create("/test/file")
	assert.NoError(t, err)

	_, err = file.Write([]byte("test data\n"))
	assert.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	buf, err := archive(fs)
	assert.NoError(t, err)

	layer, err := NewLayer(buf, "", "Test layer")
	assert.NoError(t, err)

	assert.Len(t, layer, 6656)
}
