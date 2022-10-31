package inspector

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/bisohns/saido/driver"
	log "github.com/sirupsen/logrus"
)

// DockerStatsMetrics : Metrics used by DockerStats
type DockerStatsMetrics struct {
	ContainerID   string
	ContainerName string
	CPU           float64
	MemUsage      float64
	Limit         float64
	MemPercent    float64
	Pid           int
}

// DockerStats : Parsing the `docker stats` output for container monitoring
type DockerStats struct {
	Driver  *driver.Driver
	Command string
	// We want do display disk values in GB
	DisplayByteSize string
	// Values of metrics being read
	Values []DockerStatsMetrics
}

// Parse : run custom parsing on output of the command
/*
CONTAINER  CPU %     MEM USAGE / LIMIT     MEM %     NET I/O             BLOCK I/O             PIDS
redis1     0.07%     796 KB / 64 MB        1.21%     788 B / 648 B       3.568 MB / 512 KB     2
redis2     0.07%     2.746 MB / 64 MB      4.29%     1.266 KB / 648 B    12.4 MB / 0 B         3
*/
func (i *DockerStats) Parse(output string) {
	var values []DockerStatsMetrics
	var splitChars string
	details := (*i.Driver).GetDetails()
	if details.IsWindows {
		splitChars = "\r\n"
	} else {
		splitChars = "\n"
	}
	log.Debug("Parsing output string in DockerStats inspector")
	lines := strings.Split(output, splitChars)
	for index, line := range lines {
		// skip title line
		if index == 0 {
			continue
		}
		columns := strings.Fields(line)
		if len(columns) == 14 {
			cpu, err := strconv.ParseFloat(strings.TrimSuffix(columns[2], "%"), 64)
			if err != nil {
				log.Fatal("Could not parse cpu for docker stats")
			}
			memory, err := strconv.ParseFloat(strings.TrimSuffix(columns[6], "%"), 64)
			if err != nil {
				log.Fatal("Could not parse cpu for docker stats")
			}
			col := []string{
				columns[3],
				columns[5],
			}
			pid, err := strconv.Atoi(columns[13])
			if err != nil {
				log.Debug("Could not parse pid for docker stats, probably on windows")
				pid = 0
			}
			value := i.createMetric(col, columns[0], columns[1], cpu, memory, pid)
			values = append(values, value)
		}
	}
	i.Values = values
}

func (i DockerStats) createMetric(
	columns []string,
	containerID string,
	containerName string,
	cpu float64,
	memory float64,
	pid int) DockerStatsMetrics {
	lastMem := len(columns[0]) - 3
	lastLim := len(columns[1]) - 3
	memusageSize := columns[0][lastMem:]
	limitSize := columns[1][lastLim:]
	return DockerStatsMetrics{
		ContainerID:   containerID,
		ContainerName: containerName,
		CPU:           cpu,
		MemUsage:      NewByteSize(columns[0][:lastMem], memusageSize).format(i.DisplayByteSize),
		Limit:         NewByteSize(columns[1][:lastLim], limitSize).format(i.DisplayByteSize),
		MemPercent:    memory,
		Pid:           pid,
	}
}

func (i *DockerStats) SetDriver(driver *driver.Driver) {
	i.Driver = driver
}

func (i DockerStats) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *DockerStats) Execute() ([]byte, error) {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
		return json.Marshal(i.Values)
	}
	return []byte(""), err
}

// NewDockerStats : Initialize a new DockerStats instance
func NewDockerStats(driver *driver.Driver, _ ...string) (Inspector, error) {
	var dockerstats Inspector
	details := (*driver).GetDetails()
	if !(details.IsLinux || details.IsDarwin || details.IsWindows) {
		return nil, errors.New("Cannot use LoadAvgDarwin on drivers outside (linux, darwin, windows)")
	}
	dockerstats = &DockerStats{
		Command: `docker stats --no-stream`,
	}
	dockerstats.SetDriver(driver)
	return dockerstats, nil
}
