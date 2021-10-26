package inspector

import (
	log "github.com/sirupsen/logrus"
	"strconv"
)

// ResponseTimeMetrics : Metrics used by ResponseTime
type ResponseTimeMetrics struct {
	Seconds float64
}

// ResponseTime : Parsing the `web` output for response time
type ResponseTime struct {
	fields
	// Values of metrics being read
	Values ResponseTimeMetrics
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

// NewResponseTime : Initialize a new ResponseTime instance
func NewResponseTime() *ResponseTime {
	return &ResponseTime{
		fields: fields{
			Type:    Command,
			Command: `response`,
		},
	}

}
