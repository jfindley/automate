package core

type RunControl struct {
	begin chan bool
	end   chan bool
}

func (r RunControl) Register() func() {
	// Wait for signal to start
	<-r.begin

	// Return function to signal when done
	return func() {
		close(r.end)
	}
}

func (r RunControl) Start() {
	close(r.begin)
}

func (r RunControl) Wait() {
	<-r.end
}

func NewRunControl() *RunControl {
	r := new(RunControl)
	r.begin = make(chan bool)
	r.end = make(chan bool)
	return r
}
