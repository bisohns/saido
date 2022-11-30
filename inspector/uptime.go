package inspector

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bisohns/saido/driver"
	log "github.com/sirupsen/logrus"
)

// UptimeMetrics : Metrics used by Uptime
type UptimeMetrics struct {
	Up float64
	// Idle time will not be less than uptime on
	// multiprocessor systems as the metric being
	// returned is the idle time from all processors
	// e.g 80 on an 8 processor system means each
	// processor has been idle for an average of 10 seconds
	Idle float64
	// % of time CPU has been idle
	IdlePercent float64
}

// UptimeLinux : Parsing the /proc/uptime output for uptime monitoring
type UptimeLinux struct {
	Driver   *driver.Driver
	FilePath string
	Values   *UptimeMetrics
}

type UptimeDarwin struct {
	Driver      *driver.Driver
	UpCommand   string
	IdleCommand string
	Values      *UptimeMetrics
}

type UptimeWindows struct {
	Driver    *driver.Driver
	UpCommand string
	Values    *UptimeMetrics
}

// Parse : run custom parsing on output of the command
/*
1545.95 12026.34
*/
func (i *UptimeLinux) Parse(output string) {
	fmt.Print(output)
	var err error
	log.Debug("Parsing output string in Uptime inspector")
	columns := strings.Fields(output)
	Up, err := strconv.ParseFloat(columns[0], 64)
	Idle, err := strconv.ParseFloat(columns[1], 64)
	if err != nil {
		log.Fatalf(`Error Parsing Uptime: %s `, err)
	}

	i.Values = &UptimeMetrics{
		Up:   Up,
		Idle: Idle,
	}
}

func (i *UptimeLinux) SetDriver(driver *driver.Driver) {
	details, _ := (*driver).GetDetails()
	if !details.IsLinux {
		panic("Cannot use UptimeLinux on drivers outside (linux)")
	}
	i.Driver = driver
}

func (i UptimeLinux) driverExec() driver.Command {
	return (*i.Driver).ReadFile
}

func (i *UptimeLinux) Execute() ([]byte, error) {
	output, err := i.driverExec()(i.FilePath)
	if err == nil {
		i.Parse(output)
		return json.Marshal(i.Values)
	}
	return []byte(""), err
}

// Parse : Parsing output of uptime commands on darwin
/*
1647709177
1646035560
34.96
*/
func (i *UptimeDarwin) Parse(output string) {
	fmt.Print(output)
	log.Debug("Parsing output string in UptimeDarwin inspector")
	output = strings.TrimSuffix(output, ",")
	lines := strings.Split(output, "\n")
	unixTime, err := strconv.Atoi(lines[0])
	switchedOn, err := strconv.Atoi(lines[1])
	idleTime, err := strconv.ParseFloat(lines[2], 64)
	if err != nil {
		panic("Could not parse times in UptimeDarwin")
	}

	i.Values = &UptimeMetrics{
		Up:          float64(unixTime - switchedOn),
		IdlePercent: idleTime,
	}
}

func (i *UptimeDarwin) SetDriver(driver *driver.Driver) {
	details, _ := (*driver).GetDetails()
	if !details.IsDarwin {
		panic("Cannot use UptimeDarwin on drivers outside (darwin)")
	}
	i.Driver = driver
}

func (i UptimeDarwin) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *UptimeDarwin) Execute() ([]byte, error) {
	upOutput, err := i.driverExec()(i.UpCommand)
	idleOutput, err := i.driverExec()(i.IdleCommand)
	if err == nil {
		upOutput = strings.TrimSpace(upOutput)
		upOutput = strings.TrimSuffix(upOutput, ",")
		idleOutput = strings.TrimSpace(idleOutput)
		idleOutput = strings.TrimSuffix(idleOutput, "%")
		output := fmt.Sprintf("%s\n%s", upOutput, idleOutput)
		i.Parse(output)
		return json.Marshal(i.Values)
	}
	return []byte(""), err
}

/* Parse : SystemUpTime on windows

SystemUpTime
162054

*/
func (i *UptimeWindows) Parse(output string) {
	log.Debug("Parsing output string in UptimeWindows inspector")
	output = strings.ReplaceAll(output, "\r", "")
	output = strings.ReplaceAll(output, " ", "")
	upUnformatted := strings.Split(output, "\n")[1]
	up, err := strconv.ParseFloat(upUnformatted, 64)
	if err != nil {
		panic(err)
	}
	i.Values = &UptimeMetrics{
		Up: up,
	}
}

func (i *UptimeWindows) SetDriver(driver *driver.Driver) {
	details, _ := (*driver).GetDetails()
	if !details.IsWindows {
		panic("Cannot use UptimeWindows on drivers outside (windows)")
	}
	i.Driver = driver
}

func (i UptimeWindows) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *UptimeWindows) Execute() ([]byte, error) {
	output, err := i.driverExec()(i.UpCommand)
	if err == nil {
		i.Parse(output)
		return json.Marshal(i.Values)
	}
	return []byte(""), err
}

// NewUptime : Initialize a new Uptime instance
func NewUptime(driver *driver.Driver, _ ...string) (Inspector, error) {
	var uptime Inspector
	details, err := (*driver).GetDetails()
	if err != nil {
		return nil, err
	}
	if !(details.IsDarwin || details.IsLinux || details.IsWindows) {
		return nil, errors.New("Cannot use Uptime on drivers outside (linux, darwin, windows)")
	}
	if details.IsLinux {
		uptime = &UptimeLinux{
			FilePath: `/proc/uptime`,
		}
	} else if details.IsDarwin {
		uptime = &UptimeDarwin{
			UpCommand:   `date +%s; sysctl kern.boottime | awk '{print $5}'`,
			IdleCommand: `top -l 1 | grep "CPU usage" | awk '{print $7}'`,
		}
	} else if details.IsWindows {
		uptime = &UptimeWindows{
			UpCommand: `wmic path Win32_PerfFormattedData_PerfOS_System get SystemUptime`,
		}
	}
	uptime.SetDriver(driver)
	return uptime, nil
}
