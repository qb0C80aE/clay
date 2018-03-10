package extension

// Runtime is the interface what handles the runtime environment
// * Run boots server
type Runtime interface {
	Run()
}

var registeredRuntime Runtime

// RegisterRuntime registers the runtime
func RegisterRuntime(runtime Runtime) {
	registeredRuntime = runtime
}

// GetRegisteredRuntime returns the registered runtime
func GetRegisteredRuntime() Runtime {
	return registeredRuntime
}
