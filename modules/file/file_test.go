package file

import (
	"testing"

	"github.com/jfindley/automate/core"
)

func TestFileConfigure(t *testing.T) {
	f := new(File)

	conf := core.NewConfigInput(map[string]interface{}{})

	err := f.Configure(conf)
	if err == nil {
		t.Error("No error with missing parameters")
	}

}
