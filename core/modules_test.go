package core

import (
	"fmt"
	"golang.org/x/tools/godoc/vfs"
	"testing"
)

type testModule struct {
	data map[string]interface{}
}

func (t *testModule) Name() string {
	return "Test Module"
}

func (t *testModule) Configure(in Input) error {
	t.data = in.Data()
	return nil
}

func (t *testModule) Run(r ResponseWriter) {
	r.Success(true)
	r.Changed(true)
	r.Message("info", "Run completed successfully with data:", t.data["input"])
}

func (t *testModule) RunTest(fs vfs.FileSystem, r ResponseWriter) {
	t.Run(r)
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
	m.Run(&r)

	if r.level != "info" {
		t.Error("Bad level:", r.level)
	}
	if r.message != "[Run completed successfully with data: test string]" {
		t.Error("Bad message text:", r.message)
	}
	if !r.ok || !r.callbacks {
		t.Error("Bad status:", r.ok, r.callbacks)
	}
}
