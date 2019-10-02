package wails

import (
	"sync"

	"github.com/wailsapp/wails/cmd"
	"github.com/wailsapp/wails/lib/binding"
	"github.com/wailsapp/wails/lib/event"
	"github.com/wailsapp/wails/lib/interfaces"
	"github.com/wailsapp/wails/lib/ipc"
	"github.com/wailsapp/wails/lib/logger"
	"github.com/wailsapp/wails/lib/renderer"
)

// -------------------------------- Compile time Flags ------------------------------

// BuildMode indicates what mode we are in
var BuildMode = cmd.BuildModeProd

// ----------------------------------------------------------------------------------

// App defines the main application struct
type App struct {
	config         *AppConfig                // The Application configuration object
	cli            *cmd.Cli                  // In debug mode, we have a cli
	renderer       []interfaces.Renderer     // The renderer is what we will render the app to
	logLevel       string                    // The log level of the app
	ipc            interfaces.IPCManager     // Handles the IPC calls
	log            *logger.CustomLogger      // Logger
	bindingManager interfaces.BindingManager // Handles binding of Go code to renderer
	eventManager   interfaces.EventManager   // Handles all the events
	runtime        interfaces.Runtime        // The runtime object for registered structs
}

// CreateApp creates the application window with the given configuration
// If none given, the defaults are used
func CreateApp(optionalConfig ...*AppConfig) *App {
	var userConfig *AppConfig
	if len(optionalConfig) > 0 {
		userConfig = optionalConfig[0]
	}

	result := &App{
		logLevel:       "info",
		renderer:       []interfaces.Renderer{renderer.NewWebView(), &renderer.Bridge{}},
		ipc:            ipc.NewManager(),
		bindingManager: binding.NewManager(),
		eventManager:   event.NewManager(),
		log:            logger.NewCustomLogger("App"),
	}

	appconfig, err := newConfig(userConfig)
	if err != nil {
		result.log.Fatalf("Cannot use custom HTML: %s", err.Error())
	}
	result.config = appconfig

	// Set up the CLI if not in release mode
	if BuildMode != cmd.BuildModeProd {
		result.cli = result.setupCli()
	} else {
		// Disable Inspector in release mode
		result.config.DisableInspector = true
	}

	return result
}

// Run the app
func (a *App) Run() error {
	if BuildMode != cmd.BuildModeProd {
		return a.cli.Run()
	}

	a.logLevel = "error"
	err := a.start()
	if err != nil {
		a.log.Error(err.Error())
	}
	return err
}

func (a *App) start() error {

	// Set the log level
	logger.SetLogLevel(a.logLevel)

	// Log starup
	a.log.Info("Starting")

	// Initialise the renderer
	for _, r := range a.renderer {
		err := r.Initialise(a.config, a.ipc, a.eventManager)
		if err != nil {
			return err
		}
	}

	// Start event manager and give it our renderer
	a.eventManager.Start(a.renderer)

	// Start the IPC Manager and give it the event manager and binding manager
	a.ipc.Start(a.eventManager, a.bindingManager)

	// Create the runtime
	a.runtime = NewRuntime(a.eventManager, a.renderer)

	// Start binding manager and give it our renderer
	if err := a.bindingManager.Start(a.renderer, a.runtime); err != nil {
		return err
	}

	var wg sync.WaitGroup

	// Run the renderer
	for _, r := range a.renderer {
		wg.Add(1)
		go func() {
			if err := r.Run(); err != nil {
				a.log.Fatalf("Error starting renerer: %v", err)
			}
			wg.Done()
		}()
	}

	wg.Done()

	return nil
}

// Bind allows the user to bind the given object
// with the application
func (a *App) Bind(object interface{}) {
	a.bindingManager.Bind(object)
}
