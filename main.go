package main

import ()

type Module interface {
	Name() string
	Start()
	Wait()
	Run(ResponseWriter, Input)
}

type ResponseWriter interface {
	Success(bool)
	TriggerCallbacks(bool)
	Message(string, ...interface{})
}

type Input interface {
	Data() map[string]interface{}
}

func main() {
}
