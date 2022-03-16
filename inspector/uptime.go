package inspector

import (
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

// Parse : run custom parsing on output of the command
func (i *UptimeLinux) Parse(output string) {
	var err error
	log.Debug("Parsing ouput string in Uptime inspector")
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
	details := (*driver).GetDetails()
	if !details.IsLinux {
		panic("Cannot use UptimeLinux on drivers outside (linux)")
	}
	i.Driver = driver
}

func (i UptimeLinux) driverExec() driver.Command {
	return (*i.Driver).ReadFile
}

func (i *UptimeLinux) Execute() {
	output, err := i.driverExec()(i.FilePath)
	if err == nil {
		i.Parse(output)
	}
}

func (i *UptimeDarwin) Parse(output string) {
	log.Debug("Parsing ouput string in UptimeDarwin inspector")
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
	details := (*driver).GetDetails()
	if !details.IsDarwin {
		panic("Cannot use UptimeDarwin on drivers outside (darwin)")
	}
	i.Driver = driver
}

func (i UptimeDarwin) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *UptimeDarwin) Execute() {
	upOutput, err := i.driverExec()(i.UpCommand)
	idleOutput, err := i.driverExec()(i.IdleCommand)
	if err == nil {
		upOutput = strings.TrimSuffix(upOutput, ",")
		idleOutput = strings.TrimSuffix(idleOutput, "%")
		output := fmt.Sprintf("%s\n%s", upOutput, idleOutput)
		i.Parse(output)
	}
}

//TODO: Windows equivalent of uptime

// NewUptime : Initialize a new Uptime instance
func NewUptime(driver *driver.Driver) Inspector {
	var uptime Inspector
	details := (*driver).GetDetails()
	if !(details.IsDarwin || details.IsLinux) {
		panic("Cannot use Uptime on drivers outside (linux, darwin)")
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
	}
	uptime.SetDriver(driver)
	return uptime
}
