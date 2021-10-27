package inspector

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// LoadAvgMetrics : Metrics used by LoadAvg
type LoadAvgMetrics struct {
	Load1M  float64
	Load5M  float64
	Load15M float64
}

// LoadAvg : Parsing the /proc/loadavg output for disk monitoring
type LoadAvg struct {
	fields
	Values LoadAvgMetrics
}

// Parse : run custom parsing on output of the command
func (i *LoadAvg) Parse(output string) {
	var err error
	log.Debug("Parsing ouput string in LoadAvg inspector")
	columns := strings.Fields(output)
	Load1M, err := strconv.ParseFloat(columns[0], 64)
	Load5M, err := strconv.ParseFloat(columns[1], 64)
	Load15M, err := strconv.ParseFloat(columns[2], 64)
	if err != nil {
		log.Fatalf(`Error Parsing LoadAvg: %s `, err)
	}

	i.Values = LoadAvgMetrics{
		Load1M,
		Load5M,
		Load15M,
	}
}

// NewLoadAvg : Initialize a new LoadAvg instance
func NewLoadAvg() *LoadAvg {
	return &LoadAvg{
		fields: fields{
			Type:     File,
			FilePath: `/proc/loadavg`,
		},
	}

}
