package framework

import (
	"github.com/bbredesen/go-vk"
)

var globalChannel = make(chan EventMessage, 64)

type App interface {
	// Specify the width and height of the window before opening, and the location via top and left offset.
	// For example, (800, 600, 20, 20) will position an 800 x 600 pixel window 20 pixels from the top and from the left
	// edge of the screen. Note that the renderable surface may not be the full width and height depending on
	// the underlying window system. If not called, then the default window is 640x480 at an arbitrary screen location.
	// Setting these values to a negative number will trigger the default behavior.
	SetWindowParams(width, height, left, top int)

	GetEventChannel() <-chan EventMessage

	// Run starts the OS window message loop. This needs to be called from the main thread and will block until
	// the window is closed.
	Run() error

	// OkToClose notifies the window system that it is ok to close the window at this point in time. This should be
	// called by the app once a ET_Sys_Closed message is received. The window system will cause this message to be sent
	// when the user has requested to close the window. The OkToClose callback is needed to avoid destroying the window
	// in the middle of a draw call, which will cause a crash instead of a clean shutdown.
	// The app should exit the message loop after calling this function.
	OkToClose(handle uintptr)

	GetRequiredInstanceExtensions() []string

	// Creates a surface for the window that has been opened and returns the SurfaceKHR handle. Delegates to
	// OS-specific code so that the user application can run independent from the host OS.
	DelegateCreateSurface(instance vk.Instance) (vk.SurfaceKHR, error)
}

type sharedApp struct {
	title                                string
	reqWidth, reqHeight, reqLeft, reqTop int
}

func newSharedApp(windowTitle string) sharedApp {
	return sharedApp{
		title:     windowTitle,
		reqWidth:  -1,
		reqHeight: -1,
		reqLeft:   -1,
		reqTop:    -1,
	}

}

func (app *sharedApp) GetEventChannel() <-chan EventMessage {
	return globalChannel
}

func (app *sharedApp) SetWindowParams(width, height, left, top int) {
	app.reqWidth = width
	app.reqHeight = height
	app.reqLeft = left
	app.reqTop = top
}

func (app *sharedApp) ExampleMessageLoop() {
	ch := app.GetEventChannel()
	m := <-ch // Block on the channel until the window has been created

	if m.Type != ET_Sys_Created {
		panic("expected ET_Sys_Create to start message loop")
	}

	// This is where you would create an instance, surface, set up the swapchain, etc. Your app should
	// contain all of your Vulkan handles, state, etc.
	// app.InitVulkan()

	for {
	messageLoop:
		for {
			select {
			case m = <-ch:
				switch m.Type {
				case ET_Sys_Closed:
					// This is where you would cleanup and destroy all Vulkan objects
					// app.CleanupVulkan()

					// sharedApp does not implement framework.App, but windowsApp, darwinApp, etc. do.
					// Your app should embed the framework.App interface returned from calling
					// framework.NewApp("Window Title"), which will give you access to OkToClose
					// app.OkToClose(m.SystemEvent.HandleForSurface)
					return

				case ET_Mouse_Move:
					// fmt.Println("mouse move recieved")

				}
			default: // Channel is empty
				break messageLoop

			}
		}

		// After handling all messages in the queue, you can draw a frame
		// app.drawFrame()
	}
}
