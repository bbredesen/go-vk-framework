# package framework: Cross Platform Window Interface

This package implements a simple interface for opening a window on
different OS platforms, and for receiving certain events from the user and
window system. Not all event messages have been implemented yet! 

Currently handled events:
* User events: Key press, mouse movement, mouse enter/exit
* Window events: resize, minimize, lose/gain focus, create/destroy

This package is intended to be used as a submodule in programs using `go-vk`, to quickly
get a window open and ready for Vulkan rendering. Add this repo as a submodule to your
project and then build the native delegate library with `make`:

```
mkdir framework/
cd framework/
git submodule add https://github.com/bbredesen/go-vk-framework
cd delegate/
make
```


To use this pacakge in your Go app, do the following:

```go
    app := shared.NewApp()
    
    // You have to write this, there is a sample ExampleMessageLoop in app.go 
    go myCustomEventLoop(app.GetEventChannel()) 

    app.Run() // Run() must execute on the main thread and it blocks until the window is closed.
```

## Events

Events are fed into a channel which can be read by user code. The recommended
procedure for each frame is to read all events (i.e. empty the channel)
before rendering the current frame. This framework does not attempt to bubble or
delgate events in any way.

## Message Dispatch Loop

Internally, the framework receives OS event messages and transforms them to a common format before passing to
your code:

1. Receive a message from the OS.
2. Build the internal message based on the OS-provided details.
3. Send the message to the message channel.

A final ET_Sys_Close message will be added to the channel when a request to close the window is received, after which
the channel will be closed by the sender. It is up to the programmer to clean up any Vulkan resources and call OkToClose
before exiting the event loop. The OkToClose function signals to the OS-specific delgate code that your application has 
completed cleanup and will not attempt to draw another frame, allowing the window to be safely destroyed.
