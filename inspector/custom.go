package inspector

import (
	"fmt"

	"github.com/bisohns/saido/driver"
	log "github.com/sirupsen/logrus"
)

// CustomMetrics : Metrics used by Custom
type CustomMetrics struct {
	Output string
}

// Custom : Parsing the custom command output for disk monitoring
type Custom struct {
	Driver  *driver.Driver
	Values  CustomMetrics
	Command string
}

// Parse : run custom parsing on output of the command
func (i *Custom) Parse(output string) {
	log.Debug("Parsing ouput string in Custom inspector")
	i.Values = i.createMetric(output)
}

func (i Custom) createMetric(output string) CustomMetrics {
	return CustomMetrics{
		Output: output,
	}
}

func (i *Custom) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if details.IsWeb {
		panic(fmt.Sprintf("Cannot use Custom(%s) on web", i.Command))
	}
	i.Driver = driver
}

func (i Custom) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *Custom) Execute() {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
	}
}

// NewCustom : Initialize a new Custom instance
func NewCustom(custom string, driver *driver.Driver) Inspector {
	var customInspector Inspector
	details := (*driver).GetDetails()
	if details.IsWeb {
		panic(fmt.Sprintf("Cannot use Custom(%s) on web", custom))
	}
	customInspector = &Custom{
		Command: custom,
	}
	customInspector.SetDriver(driver)
	return customInspector
}
