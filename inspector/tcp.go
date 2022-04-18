// Check if TCP ports are open or not
package inspector

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bisohns/saido/driver"
	"github.com/mum4k/termdash/widgetapi"
	"github.com/mum4k/termdash/widgets/barchart"
	log "github.com/sirupsen/logrus"
)

// TCPMetrics : Metrics obtained by tcp monitoring on darwin
type TCPMetrics struct {
	// Ports map a port to a status string
	// e.g {8081: "LISTEN"}
	Ports map[int]string
}

// TCPDarwin : TCP specfic to Darwin
type TCPDarwin struct {
	Command string
	Driver  *driver.Driver
	Values  TCPMetrics
	// FIXME: Get proper graph
	Widget *barchart.BarChart
}

// TCPLinux : TCP specfic to Linux
type TCPLinux struct {
	Command string
	Driver  *driver.Driver
	Values  TCPMetrics
	// FIXME: Get proper graph
	Widget *barchart.BarChart
}

// TCPWin : TCP specfic to Windows
type TCPWin struct {
	Command string
	Driver  *driver.Driver
	Values  TCPMetrics
	// FIXME: Get proper graph
	Widget *barchart.BarChart
}

/* Parse : parsing the following kind of output
Active Internet connections (including servers)
Proto Recv-Q Send-Q  Local Address          Foreign Address        (state)
tcp4       0      0  127.0.0.1.53300        127.0.0.1.59972        ESTABLISHED
tcp4       0      0  192.168.1.172.59964    162.247.243.147.443    SYN_SENT
tcp4       0      0  192.168.1.172.59931    13.224.227.146.443     ESTABLISHED
tcp4       0      0  127.0.0.1.59905        127.0.0.1.53300        CLOSE_WAIT
*/
func (i *TCPDarwin) Parse(output string) {
	ports := make(map[int]string)
	lines := strings.Split(output, "\n")
	for index, line := range lines {
		// skip title lines
		if index == 0 || index == 1 {
			continue
		}
		columns := strings.Fields(line)
		if len(columns) > 5 {
			status := columns[5]
			address := strings.Split(columns[3], ".")
			portString := address[len(address)-1]
			port, err := strconv.Atoi(portString)
			if err != nil {
				log.Fatal("Could not parse port number in TCPDarwin")
			}
			ports[port] = status

		}
	}
	i.Values.Ports = ports
}

func (i *TCPDarwin) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsDarwin {
		panic("Cannot use TCPDarwin on drivers outside (darwin)")
	}
	i.Driver = driver
}

func (i *TCPDarwin) GetWidget() widgetapi.Widget {
	if i.Widget == nil {
	}
	return i.Widget
}

func (i *TCPDarwin) UpdateWidget() error {
	i.Execute()
	return nil
}

func (i TCPDarwin) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *TCPDarwin) Execute() {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
	}
}

/*
Parse for output
State     Recv-Q    Send-Q       Local Address:Port     Peer Address:Port       Process
LISTEN      0         5            127.0.0.1:45481         0.0.0.0:*
LISTEN      0        4096         127.0.0.53%lo:53         0.0.0.0:*
LISTEN      0         5              127.0.0.1:631         0.0.0.0:*
ESTAB       0         0        192.168.1.106:37986      198.252.206.25:443
CLOSE-WAIT  1         0            127.0.0.1:54638         127.0.0.1:45481

*/
func (i *TCPLinux) Parse(output string) {
	ports := make(map[int]string)
	lines := strings.Split(output, "\n")
	for index, line := range lines {
		// skip title lines
		if index == 0 {
			continue
		}
		columns := strings.Fields(line)
		if len(columns) >= 5 {
			fmt.Print(columns)
			status := columns[0]
			address := strings.Split(columns[3], ":")
			portString := address[len(address)-1]
			port, err := strconv.Atoi(portString)
			if err != nil {
				log.Fatal("Could not parse port number in TCPLinux")
			}
			ports[port] = status

		}
	}
	i.Values.Ports = ports
}

func (i *TCPLinux) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsLinux {
		panic("Cannot use TCPLinux on drivers outside (linux)")
	}
	i.Driver = driver
}

func (i *TCPLinux) GetWidget() widgetapi.Widget {
	if i.Widget == nil {
	}
	return i.Widget
}

func (i *TCPLinux) UpdateWidget() error {
	i.Execute()
	return nil
}

func (i TCPLinux) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *TCPLinux) Execute() {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
	}
}

/* Parse for output

Active Connections

  Proto  Local Address          Foreign Address        State
  TCP    0.0.0.0:135            0.0.0.0:0              LISTENING
  TCP    0.0.0.0:445            0.0.0.0:0              LISTENING
  TCP    0.0.0.0:5040           0.0.0.0:0              LISTENING
  TCP    0.0.0.0:5700           0.0.0.0:0              LISTENING
  TCP    0.0.0.0:6646           0.0.0.0:0              LISTENING
  TCP    0.0.0.0:49664          0.0.0.0:0              LISTENING
*/
func (i *TCPWin) Parse(output string) {
	ports := make(map[int]string)
	lines := strings.Split(output, "\n")
	for index, line := range lines {
		// skip title lines
		if index == 0 || index == 1 || index == 3 {
			continue
		}
		columns := strings.Fields(line)
		if len(columns) > 3 {
			status := columns[3]
			address := strings.Split(columns[1], ":")
			portString := address[len(address)-1]
			port, err := strconv.Atoi(portString)
			if err != nil {
				log.Fatal("Could not parse port number in TCPWin")
			}
			ports[port] = status

		}
	}
	i.Values.Ports = ports
}

func (i *TCPWin) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsWindows {
		panic("Cannot use TCPWin on drivers outside (windows)")
	}
	i.Driver = driver
}

func (i *TCPWin) GetWidget() widgetapi.Widget {
	if i.Widget == nil {
	}
	return i.Widget
}

func (i *TCPWin) UpdateWidget() error {
	i.Execute()
	return nil
}

func (i TCPWin) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *TCPWin) Execute() {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
	}
}

// NewTCP: Initialize a new TCP instance
func NewTCP(driver *driver.Driver, _ ...string) (Inspector, error) {
	var tcp Inspector
	details := (*driver).GetDetails()
	if !(details.IsLinux || details.IsDarwin || details.IsWindows) {
		return nil, errors.New("Cannot use TCP on drivers outside (linux, darwin, windows)")
	}
	if details.IsDarwin {
		tcp = &TCPDarwin{
			Command: `netstat -anp tcp`,
		}
	} else if details.IsLinux {
		tcp = &TCPLinux{
			Command: `ss -tan`,
		}
	} else if details.IsWindows {
		tcp = &TCPWin{
			Command: `netstat -anp tcp`,
		}
	}
	tcp.SetDriver(driver)
	return tcp, nil
}
