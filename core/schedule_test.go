package core

import (
	"testing"
	"time"
)

func TestRunControl(t *testing.T) {

	ctl := NewRunControl()

	work := make(chan bool, 1)

	go func() {
		defer ctl.Register()()
		work <- true
		time.Sleep(1 * time.Second)
	}()

	// Ensure that work happens after start and before end
	startTimeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
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
		time.Sleep(1 * time.Second)
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
