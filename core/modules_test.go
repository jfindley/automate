package core

import (
	"fmt"
	"testing"
)

type testModule struct {
}

func (t testModule) Name() string {
	return "Test Module"
}

func (t testModule) Run(r ResponseWriter, i Input) {
	r.Success(true)
	r.TriggerCallbacks(true)
	r.Message("info", "Run completed successfully with data:", i.Data()["input"])
}

func (t testModule) Start() {}

func (t testModule) Wait() {}

func NewTestModule() testModule {
	return testModule{}
}

type testResponse struct {
	ok        bool
	callbacks bool
	message   string
	level     string
}

func (t *testResponse) Success(in bool) {
	t.ok = in
}

func (t *testResponse) TriggerCallbacks(in bool) {
	t.callbacks = in
}

func (t *testResponse) Message(level string, messages ...interface{}) {
	t.message = fmt.Sprint(messages)
	t.level = level
}

type testInput struct {
	data map[string]interface{}
}

func (t testInput) Data() map[string]interface{} {
	return t.data
}

// This just tests the module structure is intact
func TestModule(t *testing.T) {
	var (
		m Module
		r testResponse
		i testInput
	)

	m = NewTestModule()
	i.data = map[string]interface{}{"input": "test string"}

	m.Run(&r, i)

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
