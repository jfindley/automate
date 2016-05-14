package docker

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerConfig(t *testing.T) {
	testCfg := &ContainerConfig{
		User: "test",
	}

	v2Test := LayerConfig{
		V2ContainerConfig: testCfg,
	}

	v1Test := LayerConfig{
		V1ContainerConfig: testCfg,
	}

	assert.Equal(t, "test", v2Test.ContainerConfig().User)
	assert.Equal(t, "test", v1Test.ContainerConfig().User)
}

func TestNewID(t *testing.T) {

	id, err := newID()
	assert.NoError(t, err)

	assert.Regexp(t, regexp.MustCompile(`([0-9a-f]{64})`), id)
}
