package inspector

import (
	"fmt"

	"github.com/bisohns/saido/driver"
	"github.com/mum4k/termdash/widgetapi"
)

// memoryRender, diskRender, tcp
// Inspector : defines a particular metric supported by a driver
type Inspector interface {
	Parse(output string)
	SetDriver(driver *driver.Driver)
	Execute()
	GetWidget() widgetapi.Widget
	UpdateWidget() error
	driverExec() driver.Command
}

// NewInspector : defines the functions of each inspector to initialize
// a new instance of the inspector
type NewInspector func(driver *driver.Driver, custom ...string) (Inspector, error)

var inspectorMap = map[string]NewInspector{
	`disk`:         NewDF,
	`docker`:       NewDockerStats,
	`uptime`:       NewUptime,
	`responsetime`: NewResponseTime,
	`memory`:       NewMemInfo,
	`process`:      NewProcess,
	`custom`:       NewCustom,
	`loadavg`:      NewLoadAvg,
	`tcp`:          NewTCP,
}

// Init : initializes the specified inspector using name and driver
func Init(name string, driver *driver.Driver, custom ...string) (Inspector, error) {
	val, ok := inspectorMap[name]
	if ok {
		inspector, err := val(driver, custom...)
		if err != nil {
			return nil, err
		}
		return inspector, nil
	}
	return nil, fmt.Errorf("Cannot find inspector with name %s", name)
}

// Accepted : check if a metric is an accepted inspector
func Accepted(metric string) bool {
	allowed := false
	for key := range inspectorMap {
		if key == metric {
			allowed = true
		}
	}
	return allowed
}
