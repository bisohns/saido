package driver

import "github.com/bisoncorps/saido/inspector"

type fields struct {
	// Supported inspectors for specific driver
	Supported []inspector.Inspector
}

// Driver : specification of functions to be defined by every Driver
type Driver interface {
	readFile(path string) (string, error)
	runCommand(command string) (string, error)
	// shows the driver details, not sure if we should be showing OS name
	getDetails() string
}
