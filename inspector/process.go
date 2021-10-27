package inspector

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// ProcessMetrics : Metrics used by Process
type ProcessMetrics struct {
	Command string
	User    string
	Pid     int
	// Percentage value of CPU used
	CPU float64
	// Percentage value of memory used
	Memory float64
	// Number of seconds the process has been running
	Time int64
	TTY  string
}

// Process : Parsing the `ps -A u` output for process monitoring
type Process struct {
	fields
	// Track this particular PID
	TrackPID int
	// Values of metrics being read
	Values []ProcessMetrics
}

// Parse : run custom parsing on output of the command
func (i *Process) Parse(output string) {
	var values []ProcessMetrics
	lines := strings.Split(output, "\n")
	for index, line := range lines {
		// skip title line
		if index == 0 {
			continue
		}
		columns := strings.Fields(line)
		if len(columns) >= 10 {
			pid, err := strconv.Atoi(columns[1])
			if err != nil {
				log.Fatal("Could not parse pid in Process")
			}
			// If we are tracking only a particular ID then break loop
			if i.TrackPID != 0 && i.TrackPID == pid {
				value := i.createMetric(columns, pid)
				values = append(values, value)
				break
			} else if i.TrackPID == 0 {
				value := i.createMetric(columns, pid)
				values = append(values, value)
			}
		}
	}
	i.Values = values
}

func (i Process) createMetric(columns []string, pid int) ProcessMetrics {
	var parseErr error
	cpu, parseErr := strconv.ParseFloat(columns[2], 64)
	mem, parseErr := strconv.ParseFloat(columns[3], 64)
	unparsedTime := columns[9]
	tty := columns[6]
	minutesStr := strings.Split(unparsedTime, ":")
	minute, parseErr := strconv.Atoi(minutesStr[0])
	second, parseErr := strconv.Atoi(minutesStr[1])
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	return ProcessMetrics{
		Command: strings.Join(columns[10:], " "),
		User:    columns[0],
		CPU:     cpu,
		Memory:  mem,
		Time:    int64((minute * 60) + second),
		TTY:     tty,
	}
}

// NewProcess : Initialize a new Process instance
func NewProcess() *Process {
	return &Process{
		fields: fields{
			Type:    Command,
			Command: `ps -A u`,
		},
	}

}
