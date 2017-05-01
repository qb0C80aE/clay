package extensions

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

// RegisteredRuntime returns the registered runtime
func RegisteredRuntime() Runtime {
	return registeredRuntime
}
