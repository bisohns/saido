package inspector

import (
	"errors"
	"strconv"

	"github.com/bisohns/saido/driver"
	"github.com/mum4k/termdash/widgetapi"
	"github.com/mum4k/termdash/widgets/barchart"
	log "github.com/sirupsen/logrus"
)

// ResponseTimeMetrics : Metrics used by ResponseTime
type ResponseTimeMetrics struct {
	Seconds float64
}

// ResponseTime : Parsing the `web` output for response time
type ResponseTime struct {
	Driver  *driver.Driver
	Command string
	// Values of metrics being read
	Values ResponseTimeMetrics
	// FIXME: Get proper graph
	Widget *barchart.BarChart
}

// Parse : run custom parsing on output of the command
func (i *ResponseTime) Parse(output string) {
	log.Debug("Parsing ouput string in ResponseTime inspector")
	strconv, err := strconv.ParseFloat(output, 64)
	if err != nil {
		log.Fatal(err)
	}
	values := ResponseTimeMetrics{
		Seconds: strconv,
	}
	i.Values = values
}

func (i *ResponseTime) SetDriver(driver *driver.Driver) {
	i.Driver = driver
}

func (i *ResponseTime) GetWidget() widgetapi.Widget {
	if i.Widget == nil {
	}
	return i.Widget
}

func (i *ResponseTime) UpdateWidget() error {
	i.Execute()
	return nil
}

func (i ResponseTime) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *ResponseTime) Execute() {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
	}
}

// NewResponseTime : Initialize a new ResponseTime instance
func NewResponseTime(driver *driver.Driver, _ ...string) (Inspector, error) {
	var responsetime Inspector
	details := (*driver).GetDetails()
	if !(details.IsWeb) {
		return nil, errors.New("Cannot use response time outside driver (web)")
	}
	responsetime = &ResponseTime{
		Command: `response`,
	}
	responsetime.SetDriver(driver)
	return responsetime, nil
}
