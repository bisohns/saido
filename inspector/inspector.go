package inspector

import (
	"github.com/bisohns/saido/driver"
)

// Inspector : defines a particular metric supported by a driver
type Inspector interface {
	Parse(output string)
	SetDriver(driver *driver.Driver)
	Execute()
	driverExec() driver.Command
}
