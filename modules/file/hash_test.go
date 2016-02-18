package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var checksumTarget = []byte{
	105, 110, 105, 116, 105, 97, 108, 32, 100, 97, 116, 97,
	120, 106, 2, 247, 66, 1, 89, 3, 198, 198, 253, 133, 37,
	82, 210, 114, 145, 47, 71, 64, 225, 88, 71, 97, 138, 134,
	226, 23, 247, 31, 84, 25, 210, 94, 16, 49, 175, 238, 88,
	83, 19, 137, 100, 68, 147, 78, 176, 75, 144, 58, 104, 91,
	20, 72, 183, 85, 213, 111, 112, 26, 254, 155, 226, 206}

func TestSum(t *testing.T) {
	out := sum(testInitialData)

	assert.Equal(t, checksumTarget, out, "Checksum should match")
}
