package renderer

import (
	"fmt"

	"github.com/wailsapp/wails/lib/interfaces"
	"github.com/wailsapp/wails/lib/messages"
)

// MultiHead is a backend that hosts multiple renderers simultaneously
type MultiHead struct {
	bridge    interfaces.Renderer
	webview   interfaces.Renderer
	renderers []interfaces.Renderer
}

// Initialize creates and initializes both the wrapped webview and the bridge renderer instances.
func (d *MultiHead) Initialise(app interfaces.AppConfig, ipc interfaces.IPCManager, emgr interfaces.EventManager) error {
	// Note: WebView MUST be last
	d.renderers = []interfaces.Renderer{&Bridge{}, &WebView{}}

	for _, r := range d.renderers {
		if err := r.Initialise(app, ipc, emgr); err != nil {
			return err
		}
	}

	return nil
}

// NewMultiHead returns a new MultiHead struct
func NewMultiHead() interfaces.Renderer {
	return &MultiHead{}
}

// Run calls the Run function on each of the registered Renderers in goroutines waiting for them to finish
func (d *MultiHead) Run() error {
	errs := make([]error, len(d.renderers))

	for i, r := range d.renderers[:len(d.renderers)-1] {
		go func() {
			errs[i] = r.Run()
		}()
	}

	// The last element is the WebView which doesn't like to be placed inside another goroutine
	// which is why it is treated differently.
	if err := d.renderers[len(d.renderers)-1].Run(); err != nil {
		return err
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

// NewBinding forwards the binding request to all registered renderers
func (d *MultiHead) NewBinding(bindingName string) error {
	for _, r := range d.renderers {
		fmt.Printf("NewBinding(%s)\n", bindingName)
		if err := r.NewBinding(bindingName); err != nil {
			return err
		}
	}
	return nil
}

// Callback forwards the callback data to all registered renderers
func (d *MultiHead) Callback(data string) error {
	for _, r := range d.renderers {
		if err := r.Callback(data); err != nil {
			return err
		}
	}
	return nil
}

// NotifyEvent sends the event received to all registered renderers
func (d *MultiHead) NotifyEvent(eventData *messages.EventData) error {
	for _, r := range d.renderers {
		if err := r.NotifyEvent(eventData); err != nil {
			return err
		}
	}
	return nil
}

// SelectFile calls SelectFile on all registered Renderers
func (d *MultiHead) SelectFile() string {
	for _, r := range d.renderers {
		if val := r.SelectFile(); val != "" {
			return val
		}
	}
	return ""
}

// SelectDirectory calls SelectDirectory on all registered Renderers
func (d *MultiHead) SelectDirectory() string {
	for _, r := range d.renderers {
		if val := r.SelectDirectory(); val != "" {
			return val
		}
	}
	return ""
}

// SelectSaveFile calls SelectSaveFile on all registered Renderers
func (d *MultiHead) SelectSaveFile() string {
	for _, r := range d.renderers {
		if val := r.SelectSaveFile(); val != "" {
			return val
		}
	}
	return ""
}

// SetColour calls SetColour on all registered Renderers
func (d *MultiHead) SetColour(color string) error {
	for _, r := range d.renderers {
		if err := r.SetColour(color); err != nil {
			return err
		}
	}
	return nil
}

// EnableConsole calls EnableConsole on all registered Renderers
func (d *MultiHead) EnableConsole() {
	for _, r := range d.renderers {
		r.EnableConsole()
	}
}

// Fullscreen calls Fullscreen on all registered Renderers
func (d *MultiHead) Fullscreen() {
	for _, r := range d.renderers {
		r.Fullscreen()
	}
}

// UnFullscreen calls UnFullscreen on all registered Renderers
func (d *MultiHead) UnFullscreen() {
	for _, r := range d.renderers {
		r.UnFullscreen()
	}
}

// SetTitle calls SetTitle on all registered Renderers
func (d *MultiHead) SetTitle(title string) {
	for _, r := range d.renderers {
		r.SetTitle(title)
	}
}

// Close calls Close on all registered Renderers
func (d *MultiHead) Close() {
	for _, r := range d.renderers {
		r.Close()
	}
}
