package interfaces

// CallbackFunc is
type CallbackFunc func(string) error

// IPCManager is the event manager interface
type IPCManager interface {
	BindRenderer(Renderer)
	Dispatch(message string, f CallbackFunc)
	Start(eventManager EventManager, bindingManager BindingManager)
	Shutdown()
}
