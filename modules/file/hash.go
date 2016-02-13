package file

import (
	"github.com/codahale/blake2"
	"io/ioutil"
)

func fileChecksum(path string) ([]byte, error) {
	hash := blake2.NewBlake2B()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return hash.Sum(data), nil
}
