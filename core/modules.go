package core

type Module interface {
	Name() string
	Start()
	Wait()
	TriggeredJobs(...Module)
	Run(ResponseWriter, Input)
}

type ResponseWriter interface {
	Success(bool)
	Changed(bool)
	Message(string, ...interface{})
	Status(ModuleStatus)
}

type Input interface {
	Data() map[string]interface{}
}

type ModuleStatus struct {
	ResourceName     string
	ManagedResources []string
	ChangedResources []string
	Data             map[string]interface{}
}
