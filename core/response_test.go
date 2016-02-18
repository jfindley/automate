package core

import (
	"bytes"
	"regexp"
	"testing"

	log "github.com/Sirupsen/logrus"
)

var lvl = regexp.MustCompile(`level=(\w+)`)
var msg = regexp.MustCompile(`msg=(.+)`)

func testLevel(in string, expected string) bool {
	match := lvl.FindStringSubmatch(in)
	if len(match) < 2 || match[1] != expected {
		return false
	}
	return true
}

func testMessage(in string, expected string) bool {
	match := msg.FindStringSubmatch(in)
	// Messages will have a trailing space
	if len(match) < 2 || match[1] != expected+" " {
		println(match[1])
		return false
	}
	return true
}

func TestMessage(t *testing.T) {
	var out bytes.Buffer
    log.SetOutput(&out)
    log.SetFormatter(&log.TextFormatter{DisableColors: true})

	r := NewResponse()

	r.Success(true)
	r.Changed(true)

	if !r.Ok {
		t.Error("Bad success status")
	}

	if !r.Notify {
		t.Error("Bad notify status")
	}

	r.Message("unknown", "test")
	if !testLevel(out.String(), "warning") {
		t.Error("Bad log level")
	}
	if !testMessage(out.String(), `"Unable to parse log level: unknown"`) {
		t.Error("Bad log message")
	}

	out.Reset()
	r.Message("info", "test-info")
	if !testLevel(out.String(), "info") {
		t.Error("Bad log level")
	}
	if !testMessage(out.String(), "test-info") {
		t.Error("Bad log message")
	}

	out.Reset()
	r.Message("warn", "test-warn")
	if !testLevel(out.String(), "warning") {
		t.Error("Bad log level")
	}
	if !testMessage(out.String(), "test-warn") {
		t.Error("Bad log message")
	}

	out.Reset()
	r.Message("error", "test-error")
	if !testLevel(out.String(), "error") {
		t.Error("Bad log level")
	}
	if !testMessage(out.String(), "test-error") {
		t.Error("Bad log message")
	}

	// We don't test the higher levels as they would cause the test to exit/panic
}
