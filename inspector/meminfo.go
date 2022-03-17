package inspector

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bisohns/saido/driver"
	log "github.com/sirupsen/logrus"
)

// Metrics used by MemInfo
type MemInfoMetrics struct {
	MemTotal  float64
	MemFree   float64
	Cached    float64
	SwapTotal float64
	SwapFree  float64
}

// MemInfoLinux : Parsing the `/proc/meminfo` file output for memory monitoring
type MemInfoLinux struct {
	Driver   *driver.Driver
	FilePath string
	// The values read from the command output string are defaultly in KB
	RawByteSize string
	// We want do display disk values in GB
	DisplayByteSize string
	// Values of metrics being read
	Values *MemInfoMetrics
}

// MemInfoDarwin : Parsing `top -l 1` and `sysctl` to be able to retrieve memory details
type MemInfoDarwin struct {
	Driver         *driver.Driver
	PhysMemCommand string
	SwapCommand    string
	// The values read from the command output string are defaultly in KB
	RawByteSize string
	// We want do display disk values in GB
	DisplayByteSize string
	// Values of metrics being read
	Values *MemInfoMetrics
}

func memInfoParseOutput(output, rawByteSize, displayByteSize string) *MemInfoMetrics {
	log.Debug("Parsing ouput string in meminfo inspector")
	memTotal := getMatching("MemTotal", output)
	memFree := getMatching("MemFree", output)
	cached := getMatching("Cached", output)
	swapTotal := getMatching("SwapTotal", output)
	swapFree := getMatching("SwapFree", output)
	return createMetric(
		[]string{
			memTotal,
			memFree,
			cached,
			swapTotal,
			swapFree,
		},
		rawByteSize,
		displayByteSize,
	)
}

func getMatching(metric string, rows string) string {
	lines := strings.Split(rows, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, metric) {
			columns := strings.Fields(line)
			return columns[1]
		}
	}
	return `0`
}

func createMetric(columns []string, rawByteSize, displayByteSize string) *MemInfoMetrics {
	return &MemInfoMetrics{
		MemTotal:  NewByteSize(columns[0], rawByteSize).format(displayByteSize),
		MemFree:   NewByteSize(columns[1], rawByteSize).format(displayByteSize),
		Cached:    NewByteSize(columns[2], rawByteSize).format(displayByteSize),
		SwapTotal: NewByteSize(columns[3], rawByteSize).format(displayByteSize),
		SwapFree:  NewByteSize(columns[4], rawByteSize).format(displayByteSize),
	}
}

// Parse : run custom parsing on output of the command
func (i *MemInfoLinux) Parse(output string) {
	log.Debug("Parsing ouput string in MemInfoLinux inspector")
	i.Values = memInfoParseOutput(output, i.RawByteSize, i.DisplayByteSize)
}

func (i *MemInfoLinux) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsLinux {
		panic("Cannot use MeminfoLinux outside (linux)")
	}
	i.Driver = driver
}

func (i MemInfoLinux) driverExec() driver.Command {
	return (*i.Driver).ReadFile
}

func (i *MemInfoLinux) Execute() {
	output, err := i.driverExec()(i.FilePath)
	if err == nil {
		i.Parse(output)
	}
}

func (i *MemInfoDarwin) Parse(output string) {
	rows := strings.Split(output, "\n")
	physMemRaw := rows[0]
	swapRaw := rows[1]
	physMemCols := strings.Fields(physMemRaw)
	swapCols := strings.Fields(swapRaw)
	memUsed := strings.TrimSuffix(physMemCols[0], "M")
	memUnused := strings.TrimSuffix(physMemCols[1], "M")
	memUsedInt, err := strconv.ParseInt(memUsed, 0, 64)
	memUnusedInt, err := strconv.ParseInt(memUnused, 0, 64)
	if err != nil {
		panic("Error parsing memory on MemInfoDarwin")
	}
	memTotal := fmt.Sprintf("%d", memUsedInt+memUnusedInt)
	swapTotal := strings.TrimSuffix(swapCols[0], "M")
	swapFree := strings.TrimSuffix(swapCols[1], "M")
	//TODO: Figure out where to get cached size
	i.Values = createMetric(
		[]string{
			memTotal,
			memUnused,
			`0`,
			swapTotal,
			swapFree,
		},
		i.RawByteSize,
		i.DisplayByteSize,
	)
}

func (i *MemInfoDarwin) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsDarwin {
		panic("Cannot use MeminfoDarwin outside (darwin)")
	}
	i.Driver = driver
}

func (i MemInfoDarwin) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *MemInfoDarwin) Execute() {
	physMemOutput, err := i.driverExec()(i.PhysMemCommand)
	swapOutput, err := i.driverExec()(i.SwapCommand)

	if err == nil {
		physMemOutput = strings.TrimSuffix(physMemOutput, "\n")
		swapOutput = strings.TrimSuffix(swapOutput, "\n")
		output := fmt.Sprintf("%s\n%s", physMemOutput, swapOutput)
		i.Parse(output)
	}
}

// TODO: Windows Equivalents of MemInfo

// NewMemInfoLinux : Initialize a new MemInfoLinux instance
func NewMemInfo(driver *driver.Driver, _ ...string) (Inspector, error) {
	var meminfo Inspector
	details := (*driver).GetDetails()
	if !(details.IsLinux || details.IsDarwin) {
		return nil, errors.New("Cannot use MemInfo on drivers outside (linux, darwin)")
	}
	if details.IsLinux {
		meminfo = &MemInfoLinux{
			FilePath:        `/proc/meminfo`,
			RawByteSize:     `KB`,
			DisplayByteSize: `MB`,
		}
	} else if details.IsDarwin {
		meminfo = &MemInfoDarwin{
			PhysMemCommand:  `top -l 1 | grep PhysMem: | awk '{print $2, $6}'`,
			SwapCommand:     `sysctl -n vm.swapusage | awk '{print $3, $9}'`,
			RawByteSize:     `MB`,
			DisplayByteSize: `MB`,
		}
	}
	meminfo.SetDriver(driver)
	return meminfo, nil
}
