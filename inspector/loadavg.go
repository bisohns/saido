package inspector

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bisohns/saido/driver"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgetapi"
	"github.com/mum4k/termdash/widgets/barchart"
	log "github.com/sirupsen/logrus"
)

// LoadAvgMetrics : Metrics used by LoadAvg
type LoadAvgMetrics struct {
	Load1M  float64
	Load5M  float64
	Load15M float64
}

// LoadAvgLinux : Parsing the /proc/loadavg output for load average monitoring
type LoadAvgLinux struct {
	FilePath string
	Command  string
	Driver   *driver.Driver
	Values   *LoadAvgMetrics
	Widget   *barchart.BarChart
}

// LoadAvgDarwin : Parsing the `top` output  for Load Avg
type LoadAvgDarwin struct {
	Command string
	Driver  *driver.Driver
	Values  *LoadAvgMetrics
	Widget  *barchart.BarChart
}

// LoadAvgWin : Only grants instantaneous load metrics and not historical
type LoadAvgWin struct {
	Command string
	Driver  *driver.Driver
	Values  *LoadAvgMetrics
	Widget  *barchart.BarChart
}

func normalize(value float64, cores, max int) float64 {
	current := (value * float64(max)) / float64(cores)
	if current <= float64(max) {
		return current
	}
	return float64(max)
}

func getBarChart(size int) (*barchart.BarChart, error) {
	var labels []string
	max := 87
	min := 33
	barColors := make([]cell.Color, size)
	valueColors := make([]cell.Color, size)
	switch size {
	case 1:
		labels = []string{"%Load1M"}
	case 3:
		labels = []string{"%Load1M", "%Load5M", "%Load15M"}
	}
	for i := range valueColors {
		valueColors[i] = cell.ColorBlack
	}
	for i := range barColors {
		barColors[i] = cell.ColorNumber(rand.Intn(max-min) + min)
	}
	bc, err := barchart.New(
		barchart.BarColors(barColors),
		barchart.ValueColors(valueColors),
		barchart.ShowValues(),
		barchart.Labels(labels),
	)
	if err != nil {
		return nil, err
	}
	return bc, nil
}

func loadavgParseOutput(output string) *LoadAvgMetrics {
	var err error
	columns := strings.Fields(output)
	Load1M, err := strconv.ParseFloat(columns[0], 64)
	Load5M, err := strconv.ParseFloat(columns[1], 64)
	Load15M, err := strconv.ParseFloat(columns[2], 64)
	if err != nil {
		log.Fatalf(`Error Parsing LoadAvg: %s `, err)
	}

	return &LoadAvgMetrics{
		Load1M:  Load1M,
		Load5M:  Load5M,
		Load15M: Load15M,
	}
}

func (i *LoadAvgDarwin) GetWidget() widgetapi.Widget {
	if i.Widget == nil {
		// we need 3 bars
		i.Widget, _ = getBarChart(3)
	}
	return i.Widget
}

func (i *LoadAvgDarwin) UpdateWidget() error {
	i.Execute()
	max := 100
	values := []int{
		int(i.Values.Load1M * float64(max)),
		int(i.Values.Load5M * float64(max)),
		int(i.Values.Load15M * float64(max)),
	}
	// max value possible for a single bar is 100
	return i.Widget.Values(values, max)
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

// Parse : Parsing for darwin
/*
4.27, 5.04, 4.50
*/
func (i *LoadAvgDarwin) Parse(output string) {
	output = strings.ReplaceAll(output, ",", "")
	i.Values = loadavgParseOutput(output)
}

func (i *LoadAvgDarwin) Execute() {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
	}
}

// Parse : Linux Specific Parsing for Load Avg
/*
0.25 0.23 0.14 3/671 9362   - from /proc/loadavg
8 - from nproc
*/
func (i *LoadAvgLinux) Parse(output string) {
	splits := strings.Split(output, "\n")
	i.Values = loadavgParseOutput(splits[0])
	numberOfCores, err := strconv.Atoi(splits[1])
	max := 100
	if err == nil {
		i.Values.Load1M = normalize(i.Values.Load1M, numberOfCores, max)
		i.Values.Load5M = normalize(i.Values.Load5M, numberOfCores, max)
		i.Values.Load15M = normalize(i.Values.Load15M, numberOfCores, max)
	}
}

func (i *LoadAvgLinux) GetWidget() widgetapi.Widget {
	if i.Widget == nil {
		// we need 3 bars
		i.Widget, _ = getBarChart(3)
	}
	return i.Widget
}

func (i *LoadAvgLinux) UpdateWidget() error {
	i.Execute()
	max := 100
	values := []int{
		int(i.Values.Load1M),
		int(i.Values.Load5M),
		int(i.Values.Load5M),
	}
	// max value possible for a single bar is 100
	return i.Widget.Values(values, max)
}

func (i *LoadAvgLinux) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsLinux {
		panic("Cannot use LoadAvg on drivers outside (linux)")
	}
	i.Driver = driver
}

func (i LoadAvgLinux) driverExec() driver.Command {
	return (*i.Driver).ReadFile
}

func (i *LoadAvgLinux) Execute() {
	nprocOutput, nErr := (*i.Driver).RunCommand(i.Command)
	output, fErr := i.driverExec()(i.FilePath)
	if nErr == nil && fErr == nil {
		finalOutput := fmt.Sprintf("%s%s", output, nprocOutput)
		i.Parse(finalOutput)
	}
}

// Parse : Windows specific parsing for Windows
/*

 */
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

func (i *LoadAvgWin) GetWidget() widgetapi.Widget {
	if i.Widget == nil {
		// we need 1
		i.Widget, _ = getBarChart(1)
	}
	return i.Widget
}

func (i *LoadAvgWin) UpdateWidget() error {
	i.Execute()
	max := 100
	values := []int{
		int(i.Values.Load1M),
	}
	// max value possible for a single bar is 100
	return i.Widget.Values(values, max)
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
		loadavg = &LoadAvgLinux{
			FilePath: `/proc/loadavg`,
			Command:  `nproc`,
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
