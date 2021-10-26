package inspector

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

// Metrics used by MemInfo
type MemInfoMetrics struct {
	MemTotal  float64
	MemFree   float64
	Cached    float64
	SwapTotal float64
	SwapFree  float64
}

// MemInfo : Parsing the `/proc/meminfo` file output for memory monitoring
type MemInfo struct {
	fields
	// The values read from the command output string are defaultly in KB
	RawByteSize string
	// We want do display disk values in GB
	DisplayByteSize string
	// Values of metrics being read
	Values MemInfoMetrics
}

// Parse : run custom parsing on output of the command
func (i *MemInfo) Parse(output string) {
	log.Debug("Parsing ouput string in MemInfo inspector")
	memTotal := i.getMatching("MemTotal", output)
	memFree := i.getMatching("MemFree", output)
	cached := i.getMatching("Cached", output)
	swapTotal := i.getMatching("SwapTotal", output)
	swapFree := i.getMatching("SwapFree", output)
	i.Values = i.createMetric([]string{memTotal, memFree, cached, swapTotal, swapFree})
}

func (i MemInfo) getMatching(metric string, rows string) string {
	lines := strings.Split(rows, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, metric) {
			columns := strings.Fields(line)
			return columns[1]
		}
	}
	return `0`
}

func (i MemInfo) createMetric(columns []string) MemInfoMetrics {
	return MemInfoMetrics{
		MemTotal:  NewByteSize(columns[0], i.RawByteSize).format(i.DisplayByteSize),
		MemFree:   NewByteSize(columns[1], i.RawByteSize).format(i.DisplayByteSize),
		Cached:    NewByteSize(columns[2], i.RawByteSize).format(i.DisplayByteSize),
		SwapTotal: NewByteSize(columns[3], i.RawByteSize).format(i.DisplayByteSize),
		SwapFree:  NewByteSize(columns[4], i.RawByteSize).format(i.DisplayByteSize),
	}
}

// NewMemInfo : Initialize a new MemInfo instance
func NewMemInfo() *MemInfo {
	return &MemInfo{
		fields: fields{
			Type:     File,
			FilePath: `/proc/meminfo`,
		},
		RawByteSize:     `KB`,
		DisplayByteSize: `MB`,
	}

}
