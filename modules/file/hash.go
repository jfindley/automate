package file

import (
	"github.com/codahale/blake2"
)

func sum(data []byte) []byte {
    hash := blake2.NewBlake2B()
    return hash.Sum(data)
}