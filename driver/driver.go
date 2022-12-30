package driver

import "github.com/bisohns/saido/config"

// SystemInfo gives more insight into system details
type SystemDetails struct {
	IsWindows bool
	IsLinux   bool
	IsDarwin  bool
	IsWeb     bool
	Name      string
	Extra     string
}

type driverBase struct {
	// Polling interval between retrievals
	PollInterval int64
	Info         *SystemDetails
}

// Command represents the two commands ReadFile & RunCommand
type Command func(string) (string, error)

// Driver : specification of functions to be defined by every Driver
type Driver interface {
	ReadFile(path string) (string, error)
	RunCommand(command string) (string, error)
	// shows the driver details, not sure if we should be showing OS name
	GetDetails() (SystemDetails, error)
}

func ToDriver(conn config.Connection) Driver {
	switch conn.Type {
	case "ssh":
		return &SSH{
			User:            conn.Username,
			Host:            conn.Host,
			Port:            int(conn.Port),
			KeyFile:         conn.PrivateKeyPath,
			KeyPass:         conn.PrivateKeyPassPhrase,
			Password:        conn.Password,
			CheckKnownHosts: false,
		}
	default:
		return &Local{}
	}
}
