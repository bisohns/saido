package inspector

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type DFMetrics struct {
	size        float64
	used        float64
	available   float64
	percentFull int
}

// DF : Parsing the `df` output for memory monitoring
type DF struct {
	fields
	// The values read from the command output string are defaultly in KB
	RawByteSize string
	// We want do display disk values in GB
	DisplayByteSize string
	// Parse only device that start with this e.g /dev/sd
	DeviceStartsWith string
	// Mount point to examine
	MountPoint string
	// Values of metrics being read
	Values []DFMetrics
}

// Parse : run custom parsing on output of the command
func (i *DF) Parse(output string) {
	var values []DFMetrics
	log.Debug("Parsing ouput string in DF inspector")
	lines := strings.Split(output, "\n")
	for index, line := range lines {
		// skip title line
		if index == 0 {
			continue
		}
		columns := strings.Fields(line)
		if len(columns) == 6 {
			percent := columns[4]
			if len(percent) > 1 {
				percent = percent[:len(percent)-1]
			} else if percent == `-` {
				percent = `0`
			}
			percentInt, err := strconv.Atoi(percent)
			if err != nil {
				log.Fatalf(`Error Parsing Percent Full: %s `, err)
			}
			if columns[5] == i.MountPoint {
				values = append(values, i.createMetric(columns, percentInt))
			} else if strings.HasPrefix(columns[0], i.DeviceStartsWith) &&
				i.MountPoint == "" {
				values = append(values, i.createMetric(columns, percentInt))
			}
		}
	}
	i.Values = values
}

func (i DF) createMetric(columns []string, percent int) DFMetrics {
	return DFMetrics{
		size:        NewByteSize(columns[1], i.RawByteSize).format(i.DisplayByteSize),
		used:        NewByteSize(columns[2], i.RawByteSize).format(i.DisplayByteSize),
		available:   NewByteSize(columns[3], i.RawByteSize).format(i.DisplayByteSize),
		percentFull: percent,
	}
}

// NewDF : Initialize a new DF instance
func NewDF() *DF {
	return &DF{
		fields: fields{
			Type:    Command,
			Command: `df -a`,
		},
		RawByteSize:     `KB`,
		DisplayByteSize: `GB`,
		MountPoint:      `/`,
	}

}
