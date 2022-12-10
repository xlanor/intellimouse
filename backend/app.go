package backend

import (
	"context"
	"fmt"
	"github.com/xlanor/intellimouse/api"
	"github.com/xlanor/intellimouse/backend/internal"
	"github.com/sirupsen/logrus"
)

// App struct
type App struct {
	ctx context.Context
	Log *logrus.Logger
	Driver *api.Driver
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) Startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	a.Log = internal.NewLogger()
	a.LoadDevicesPolling()
	return
}

// domReady is called after the front-end dom has been loaded
func (a *App) DomReady(ctx context.Context) {
	// Add your action here
}

func (a *App) LoadDevices(ctx context.Context) {

	devices := api.GetDeviceInfo(api.HidStruct{})
	fmt.Println(len(devices))
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
}
