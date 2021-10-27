package inspector

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// UptimeMetrics : Metrics used by Uptime
type UptimeMetrics struct {
	Up   float64
	Idle float64
}

// Uptime : Parsing the /proc/uptime output for uptime monitoring
type Uptime struct {
	fields
	Values UptimeMetrics
}

// Parse : run custom parsing on output of the command
func (i *Uptime) Parse(output string) {
	var err error
	log.Debug("Parsing ouput string in Uptime inspector")
	columns := strings.Fields(output)
	Up, err := strconv.ParseFloat(columns[0], 64)
	Idle, err := strconv.ParseFloat(columns[1], 64)
	if err != nil {
		log.Fatalf(`Error Parsing Uptime: %s `, err)
	}

	i.Values = UptimeMetrics{
		Up,
		Idle,
	}
}

// NewUptime : Initialize a new Uptime instance
func NewUptime() *Uptime {
	return &Uptime{
		fields: fields{
			Type:     File,
			FilePath: `/proc/uptime`,
		},
	}

}
