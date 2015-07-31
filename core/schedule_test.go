package core

import (
	"testing"
	"time"
)

// Timer to prevent races during schedule test.
// Deliberately set high to allow for old hardware.
var testRaceMs = 300 * time.Millisecond

func TestRunControl(t *testing.T) {

	ctl := NewRunControl()

	work := make(chan bool, 1)

	go func() {
		defer ctl.Register()()
		work <- true
		time.Sleep(testRaceMs)
	}()

	// Ensure that work happens after start and before end
	startTimeout := make(chan bool, 1)
	go func() {
		time.Sleep(testRaceMs)
		startTimeout <- true
	}()

	select {
	case <-work:
		t.Fatal("Work was done before start was called")
	case <-startTimeout:
		break
	}

	ctl.Start()

	if _, ok := <-ctl.end; ok {
		t.Fatal("Ended before work was done")
	}

	endTimeout := make(chan bool, 1)
	go func() {
		time.Sleep(testRaceMs)
		endTimeout <- true
	}()

	select {
	case <-work:
		break
	case <-endTimeout:
		t.Fatal("Timed out while waiting for work to be done")
	}

	ctl.Wait()
}
