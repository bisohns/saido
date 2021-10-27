package inspector

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// TasklistMetrics : Metrics used by Tasklist
type TasklistMetrics struct {
	Command string
	Session string
	Pid     int
	// value of memory used
	Memory float64
}

// Tasklist : Parsing the `tasklist` output for process monitoring on windows
type Tasklist struct {
	fields
	// Track this particular PID
	TrackPID int
	// We want do display memory values in KB
	DisplayByteSize string
	// Values of metrics being read
	Values []TasklistMetrics
}

// Parse : run custom parsing on output of the command
func (i *Tasklist) Parse(output string) {
	var values []TasklistMetrics
	lines := strings.Split(output, "\r\n")
	for index, line := range lines {
		// skip title line
		if index == 0 || index == 1 {
			continue
		}
		columns := strings.Fields(line)
		lenCol := len(columns)
		if lenCol >= 6 {
			pid, err := strconv.Atoi(columns[lenCol-5])
			if err != nil {
				log.Fatal("Could not parse pid in Tasklist")
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

func (i Tasklist) createMetric(columns []string, pid int) TasklistMetrics {
	lenCol := len(columns)

	return TasklistMetrics{
		Command: strings.Join(columns[:lenCol-5], " "),
		Memory:  NewByteSize(columns[lenCol-2], `KB`).format(i.DisplayByteSize),
		Session: columns[lenCol-4],
		Pid:     pid,
	}
}

// NewTasklist : Initialize a new Tasklist instance
func NewTasklist() *Tasklist {
	return &Tasklist{
		fields: fields{
			Type:    Command,
			Command: `tasklist`,
		},
		DisplayByteSize: `KB`,
	}

}
