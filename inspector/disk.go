package inspector

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bisohns/saido/driver"
	log "github.com/sirupsen/logrus"
)

// DFMetrics : Metrics used by DF
type DFMetrics struct {
	FileSystem  string
	Size        float64
	Used        float64
	Available   float64
	PercentFull int
	// Optional Volume Name that may be available on Windows
	VolumeName string
}

// DF : Parsing the `df` output for disk monitoring
type DF struct {
	Driver  *driver.Driver
	Command string
	// The values read from the command output string are defaultly in KB
	RawByteSize string
	// We want do display disk values in GB
	DisplayByteSize string
	// Parse only device that start with this e.g /dev/sd
	DeviceStartsWith string
	// Values of metrics being read
	Values []DFMetrics
}

// Parse : run custom parsing on output of the command
/*
For Darwin it looks something like

 FileSystem    1024-blocks      Used Available Capacity iused      ifree %iused  Mounted on
 /dev/disk1s5    244679060  10984568  47579472    19%  488275 2446302325    0%   /
 devfs                 220       220         0   100%     774          0  100%   /dev
 /dev/disk1s1    244679060 179090752  47579472    80% 1205263 2445585337    0%   /System/Volumes/Data
 /dev/disk1s4    244679060   6292564  47579472    12%       7 2446790593    0%   /private/var/vm
 map auto_home           0         0         0   100%       0          0  100%   /System/Volumes/Data/home

*/
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
		if len(columns) >= 6 {
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
			if strings.HasPrefix(columns[0], i.DeviceStartsWith) {
				values = append(values, i.createMetric(columns, percentInt))
			} else {
				values = append(values, i.createMetric(columns, percentInt))
			}
		}
	}
	i.Values = values
}

func (i DF) createMetric(columns []string, percent int) DFMetrics {
	return DFMetrics{
		FileSystem:  columns[0],
		Size:        NewByteSize(columns[1], i.RawByteSize).format(i.DisplayByteSize),
		Used:        NewByteSize(columns[2], i.RawByteSize).format(i.DisplayByteSize),
		Available:   NewByteSize(columns[3], i.RawByteSize).format(i.DisplayByteSize),
		PercentFull: percent,
	}
}

func (i *DF) SetDriver(driver *driver.Driver) {
	i.Driver = driver
}

func (i DF) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *DF) Execute() {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
	}
}

// DFWin: parse `wmic logicaldisk` to satisfy Inspector interface
type DFWin struct {
	Driver  *driver.Driver
	Command string
	// The values read from the command output string are defaultly in KB
	RawByteSize string
	// We want do display disk values in GB
	DisplayByteSize string
	// Parse only device that start with this e.g /dev/sd
	DeviceStartsWith string
	// Values of metrics being read
	Values []DFMetrics
}

/* Parse : For the following windows output

Node,DeviceID,DriveType,FreeSpace,ProviderName,Size,VolumeName
IMANI,C:,3,191980253184,,288303964160,OS
*/
func (i *DFWin) Parse(output string) {
	var values []DFMetrics
	log.Debug("Parsing ouput string in DF inspector")
	lineChar := "\r"
	output = strings.TrimPrefix(output, lineChar)
	output = strings.TrimSuffix(output, lineChar)
	lines := strings.Split(output, lineChar)
	for index, line := range lines {
		// skip title line
		if index == 0 || index == 1 {
			continue
		}
		columns := strings.Split(line, ",")
		if len(columns) >= 7 {
			available, err := strconv.Atoi(columns[3])
			size, err := strconv.Atoi(columns[5])
			if err != nil {
				panic("Could not parse sizes for DFWin")
			}
			used := size - available
			percentInt := int((float64(used) / float64(size)) * 100)
			cols := []string{
				columns[1],
				fmt.Sprintf("%d", size),
				fmt.Sprintf("%d", used),
				fmt.Sprintf("%d", available),
				columns[6],
			}
			if strings.HasPrefix(columns[1], i.DeviceStartsWith) {
				values = append(values, i.createMetric(cols, percentInt))
			} else {
				values = append(values, i.createMetric(cols, percentInt))
			}
		}
	}
	i.Values = values
}

func (i DFWin) createMetric(columns []string, percent int) DFMetrics {
	return DFMetrics{
		FileSystem:  columns[0],
		Size:        NewByteSize(columns[1], i.RawByteSize).format(i.DisplayByteSize),
		Used:        NewByteSize(columns[2], i.RawByteSize).format(i.DisplayByteSize),
		Available:   NewByteSize(columns[3], i.RawByteSize).format(i.DisplayByteSize),
		VolumeName:  columns[4],
		PercentFull: percent,
	}
}

func (i *DFWin) SetDriver(driver *driver.Driver) {
	i.Driver = driver
}

func (i DFWin) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *DFWin) Execute() {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
	}
}

// NewDF : Initialize a new DF instance
func NewDF(driver *driver.Driver, _ ...string) (Inspector, error) {
	var df Inspector
	details := (*driver).GetDetails()
	if !(details.IsLinux || details.IsDarwin || details.IsWindows) {
		return nil, errors.New("Cannot use 'df' command on drivers outside (linux, darwin, windows)")
	}
	if details.IsLinux || details.IsDarwin {
		df = &DF{
			// Using -k to ensure size is
			// always reported in posix standard of 1K-blocks
			Command:         `df -a -k`,
			RawByteSize:     `KB`,
			DisplayByteSize: `GB`,
		}
	} else {
		df = &DFWin{
			// Using format to account for weird spacing
			// issues that arise on windows
			Command:         `wmic logicaldisk list brief /format:csv`,
			RawByteSize:     `B`,
			DisplayByteSize: `GB`,
		}
	}
	df.SetDriver(driver)
	return df, nil
}
