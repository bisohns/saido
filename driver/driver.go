package driver

type fields struct {
	// Supported inspector representations for specific driver
	Supported []string
	// Selected inspector representations
	Selected []string
	// Polling interval between retrievals
	PollInterval int64
}

// Driver : specification of functions to be defined by every Driver
type Driver interface {
	ReadFile(path string) (string, error)
	RunCommand(command string) (string, error)
	// shows the driver details, not sure if we should be showing OS name
	GetDetails() string
}
