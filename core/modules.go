package core

import (
	"sourcegraph.com/sourcegraph/rwvfs"
)

type Module interface {
	Name() string
	Configure(Input) error
	Run(ResponseWriter)
	RunTest(rwvfs.FileSystem, ResponseWriter)
}

type ResponseWriter interface {
	Success(bool)
	Changed(bool)
	Message(string, ...interface{})
	TriggeredJobs(...Module)
	Status(ModuleStatus)
}

type Input interface {
	Data() map[string]interface{}
	Validate(InputSchema) error
}

type ModuleStatus struct {
	ResourceName     string
	ManagedResources []string
	ChangedResources []string
	Metadata         map[string]interface{}
}
