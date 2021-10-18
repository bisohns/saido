package driver

// Driver : specification of functions to be defined by every Driver
type Driver interface {
	readFile(path string) (string, error)
	runCommand(command string) (string, error)
	// shows the driver details, not sure if we should be showing OS name
	getDetails() string
}
