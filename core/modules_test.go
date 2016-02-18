package core

import (
	"fmt"
	"sourcegraph.com/sourcegraph/rwvfs"
	"testing"
	"github.com/stretchr/testify/assert"
)

type testModule struct {
	data map[string]interface{}
}

func (t *testModule) Name() string {
	return "Test Module"
}

func (t *testModule) Configure(in Input) error {
	val, err := in.Data("input")
    t.data = make(map[string]interface{})
    t.data["input"] = val
	return err
}

func (t *testModule) Run(fs rwvfs.FileSystem, r ResponseWriter) {
	r.Success(true)
	r.Changed(true)
	r.Message("info", "Run completed successfully with data:", t.data["input"])
}

func NewTestModule() *testModule {
	return new(testModule)
}

type testResponse struct {
	ok        bool
	callbacks bool
	message   string
	level     string
	status    ModuleStatus
}

func (t *testResponse) Success(in bool) {
	t.ok = in
}

func (t *testResponse) Changed(in bool) {
	t.callbacks = in
}

func (t *testResponse) Message(level string, messages ...interface{}) {
	t.message = fmt.Sprint(messages)
	t.level = level
}

func (t *testResponse) Status(m ModuleStatus) {
	t.status = m
}

func (t *testResponse) TriggeredJobs(jobs ...Module) {

}

type testInput struct {
	data map[string]interface{}
}

func (t *testInput) Data() map[string]interface{} {
	return t.data
}

func (t *testInput) Validate(in InputSchema) error {
	return nil
}

// This just tests the module structure is intact
func TestModule(t *testing.T) {
	var (
		m Module
		r testResponse
	)

	m = NewTestModule()
	in := NewConfigInput(map[string]interface{}{"input": "test string"})

	m.Configure(in)
    
    fs := rwvfs.Map(map[string]string{})
	m.Run(fs, &r)
    
    assert.Equal(t, r.level, "info")
    assert.Equal(t, "[Run completed successfully with data: test string]", r.message)
    assert.True(t, r.ok)
    assert.True(t, r.callbacks)
}
