package inspector

import (
	log "github.com/sirupsen/logrus"
)

// Mode : This specifies whether an Inspector is a command or a file
type Mode int

const (
	// Command : Inspector is a command to be executes
	Command Mode = iota
	// File : Inspector is a file to be read
	File
)

var inspectorMap = map[string]Inspector{
	`disk`:    NewDF(),
	`meminfo`: NewMemInfo(),
}

type fields struct {
	// Specify a mode for the Inspector
	Type Mode
	// File path to read
	FilePath string
	// Command to execute
	Command string
}

func (f *fields) String() string {
	value := `None`
	if f.Type == Command {
		value = f.Command
	} else if f.Type == File {
		value = f.FilePath
	}
	return value
}

// Inspector : defines a particular metric supported by a driver
type Inspector interface {
	Parse(output string)
}

// GetInspector : obtain an initialized inspector using name
func GetInspector(name string) Inspector {
	val, ok := inspectorMap[name]
	if !ok {
		log.Fatalf(`%s inspector not found`, name)
	}
	return val
}
