package inspector

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bisohns/saido/driver"
	log "github.com/sirupsen/logrus"
)

// LoadAvgMetrics : Metrics used by LoadAvg
type LoadAvgMetrics struct {
	Load1M  float64
	Load5M  float64
	Load15M float64
}

// LoadAvg : Parsing the /proc/loadavg output for load average monitoring
type LoadAvg struct {
	FilePath string
	Driver   *driver.Driver
	Values   *LoadAvgMetrics
}

type LoadAvgDarwin struct {
	Command string
	Driver  *driver.Driver
	Values  *LoadAvgMetrics
}

// LoadAvgWin : Only grants instantaneous load metrics and not historical
type LoadAvgWin struct {
	Command string
	Driver  *driver.Driver
	Values  *LoadAvgMetrics
}

func loadavgParseOutput(output string) *LoadAvgMetrics {
	var err error
	log.Debug("Parsing ouput string in LoadAvg inspector")
	columns := strings.Fields(output)
	Load1M, err := strconv.ParseFloat(columns[0], 64)
	Load5M, err := strconv.ParseFloat(columns[1], 64)
	Load15M, err := strconv.ParseFloat(columns[2], 64)
	if err != nil {
		log.Fatalf(`Error Parsing LoadAvg: %s `, err)
	}

	return &LoadAvgMetrics{
		Load1M,
		Load5M,
		Load15M,
	}
}

func (i *LoadAvgDarwin) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsDarwin {
		panic("Cannot use LoadAvgDarwin on drivers outside (darwin)")
	}
	i.Driver = driver
}

func (i LoadAvgDarwin) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *LoadAvgDarwin) Parse(output string) {
	output = strings.TrimSuffix(output, "}")
	output = strings.TrimPrefix(output, "{")
	i.Values = loadavgParseOutput(output)
}

func (i *LoadAvgDarwin) Execute() {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
	}
}

func (i *LoadAvg) Parse(output string) {
	i.Values = loadavgParseOutput(output)
}

func (i *LoadAvg) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsLinux {
		panic("Cannot use LoadAvg on drivers outside (linux)")
	}
	i.Driver = driver
}

func (i LoadAvg) driverExec() driver.Command {
	return (*i.Driver).ReadFile
}

func (i *LoadAvg) Execute() {
	output, err := i.driverExec()(i.FilePath)
	if err == nil {
		i.Parse(output)
	}
}

func (i *LoadAvgWin) Parse(output string) {
	output = strings.ReplaceAll(output, "\r", "")
	output = strings.ReplaceAll(output, " ", "")
	columns := strings.Split(output, "\n")
	// Only instantaneous metrics available so append the
	// rest as zero
	output = columns[1]
	output = fmt.Sprintf("%s 0 0", output)
	i.Values = loadavgParseOutput(output)
}

func (i *LoadAvgWin) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsWindows {
		panic("Cannot use LoadAvgWin on drivers outside (windows)")
	}
	i.Driver = driver
}

func (i LoadAvgWin) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *LoadAvgWin) Execute() {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
	}
}

// NewLoadAvg : Initialize a new LoadAvg instance
func NewLoadAvg(driver *driver.Driver, _ ...string) (Inspector, error) {
	var loadavg Inspector
	details := (*driver).GetDetails()
	if !(details.IsLinux || details.IsDarwin || details.IsWindows) {
		return nil, errors.New("Cannot use LoadAvg on drivers outside (linux, darwin)")
	}
	if details.IsLinux {
		loadavg = &LoadAvg{
			FilePath: `/proc/loadavg`,
		}
	} else if details.IsDarwin {
		loadavg = &LoadAvgDarwin{
			//      Command: `sysctl -n vm.loadavg | awk '{ printf "%.2f %.2f %.2f ", $2, $3, $4 }'`,
			Command: `top -l 1 | grep "Load Avg:" | awk '{print $3, $4, $5}'`,
		}
	} else if details.IsWindows {
		loadavg = &LoadAvgWin{
			Command: `wmic cpu get loadpercentage`,
		}
	}
	loadavg.SetDriver(driver)
	return loadavg, nil

}
