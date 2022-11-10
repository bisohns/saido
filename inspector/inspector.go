package inspector

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bisohns/saido/driver"
)

// CustomCommand : every custom command must be prefixed by this
var CustomCommand = `custom`

// Inspector : defines a particular metric supported by a driver
type Inspector interface {
	Parse(output string)
	SetDriver(driver *driver.Driver)
	Execute() ([]byte, error)
	driverExec() driver.Command
}

type NewInspector func(driver *driver.Driver, custom ...string) (Inspector, error)

var inspectorMap = map[string]NewInspector{
	`disk`:        NewDF,
	`docker`:      NewDockerStats,
	`uptime`:      NewUptime,
	`memory`:      NewMemInfo,
	`process`:     NewProcess,
	`loadavg`:     NewLoadAvg,
	`tcp`:         NewTcp,
	CustomCommand: NewCustom,
	// NOTE: Inactive for now
	`responsetime`: NewResponseTime,
}

// Valid : checks if inspector is a valid inspector
func Valid(name string) bool {
	for key := range inspectorMap {
		if name == key || strings.HasPrefix(name, CustomCommand) {
			return true
		}
	}
	return false
}

// Init : initializes the specified inspector using name and driver
func Init(name string, driver *driver.Driver, custom ...string) (Inspector, error) {
	if strings.HasPrefix(name, CustomCommand) {
		name = "custom"
	}
	val, ok := inspectorMap[name]
	if ok {
		inspector, err := val(driver, custom...)
		if err != nil {
			return nil, err
		}
		return inspector, nil
	}
	return nil, errors.New(fmt.Sprintf("Cannot find inspector with name %s", name))
}
