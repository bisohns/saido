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

type ProcessMetricsWin struct {
	Command     string
	SessionName string
	Pid         int
	Memory      float64
}

// Process : Parsing the `ps axu` output for process monitoring
type Process struct {
	Driver  *driver.Driver
	Command string
	// Track this particular PID
	TrackPID int
	// Values of metrics being read
	Values []ProcessMetrics
}

// ProcessWin : Parsing the `tasklist` output for process monitoring on windows
type ProcessWin struct {
	Driver   *driver.Driver
	Command  string
	TrackPID int
	Values   []ProcessMetricsWin
}

func (i *Process) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !(details.IsLinux || details.IsDarwin) {
		panic("Cannot use Process on drivers outside (linux, darwin)")
	}
	i.Driver = driver
}

func (i Process) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *Process) Execute() ([]byte, error) {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
		return json.Marshal(i.Values)
	}
	return []byte(""), err
}

// Parse : run custom parsing on output of the command
/*
for darwin it looks something like

USER               PID  %CPU %MEM      VSZ    RSS   TT  STAT STARTED      TIME COMMAND
formplus           980 111.4  0.1 12845480   5132   ??  S    28Feb22 4196:13.21 com.docker.hyperkit -A -u -F vms/0/hyperkit.pid -c
formplus           864  10.2  0.2  5066712  20516   ??  S    28Feb22 419:29.39 /Applications/Docker.app/Contents/MacOS/com.docker
formplus         83064   1.6  0.1  4301640   5448 s008  Ss   Sat09PM   0:43.12 -zsh


for linux it looks like

USER         PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root           1  0.0  0.0 167580 11440 ?        Ss   18:07   0:00 /sbin/init splash
root           2  0.0  0.0      0     0 ?        S    18:07   0:00 [kthreadd]
root           3  0.0  0.0      0     0 ?        I<   18:07   0:00 [rcu_gp]
*/
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
	second, parseErr := strconv.ParseFloat(minutesStr[1], 64)
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	return ProcessMetrics{
		Command: strings.Join(columns[10:], " "),
		User:    columns[0],
		CPU:     cpu,
		Memory:  mem,
		Time:    int64((minute * 60) + int(second)),
		TTY:     tty,
	}
}

/* Parse for the following

Image Name                     PID Session Name        Session#    Mem Usage
========================= ======== ================ =========== ============
System Idle Process              0 Services                   0          8 K
System                           4 Services                   0      7,304 K
Secure System                  104 Services                   0     40,344 K
Registry                       168 Services                   0     93,448 K
smss.exe                       604 Services                   0      1,080 K
csrss.exe                      968 Services                   0      4,916 K
*/
func (i *ProcessWin) Parse(output string) {
	fmt.Print(output)
	var values []ProcessMetricsWin
	lines := strings.Split(output, "\r\n")
	for index, line := range lines {
		// skip title lines and ===== line
		if index == 0 || index == 1 || index == 2 {
			continue
		}
		columns := strings.Fields(line)
		colLength := len(columns)
		if colLength >= 6 {
			pidRaw := columns[colLength-5]
			pid, err := strconv.Atoi(pidRaw)
			if err != nil {
				panic("Could not parse pid for row")
			}
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

func (i *ProcessWin) createMetric(columns []string, pid int) ProcessMetricsWin {
	colLength := len(columns)
	memoryRaw := strings.Replace(columns[colLength-2], ",", "", -1)
	memory, err := strconv.ParseFloat(memoryRaw, 64)
	if err != nil {
		panic("Error parsing memory in ProcessWin")
	}
	sessionName := columns[colLength-4]
	command := strings.Join(columns[:colLength-5], " ")

	return ProcessMetricsWin{
		Command:     command,
		Pid:         pid,
		SessionName: sessionName,
		Memory:      memory,
	}
}

func (i *ProcessWin) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsWindows {
		panic("Cannot use ProcessWin on drivers outside (windows)")
	}
	i.Driver = driver
}

func (i ProcessWin) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *ProcessWin) Execute() ([]byte, error) {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
		return json.Marshal(i.Values)
	}
	return []byte(""), err
}

// NewProcess : Initialize a new Process instance
func NewProcess(driver *driver.Driver, _ ...string) (Inspector, error) {
	var process Inspector
	details := (*driver).GetDetails()
	if !(details.IsLinux || details.IsDarwin || details.IsWindows) {
		return nil, errors.New("Cannot use Process on drivers outside (linux, darwin, windows)")
	}
	if details.IsLinux || details.IsDarwin {
		process = &Process{
			Command: `ps axu`,
		}
	} else {
		process = &ProcessWin{
			Command: `tasklist`,
		}
	}
	process.SetDriver(driver)
	return process, nil
}
