// Check if TCP ports are open or not
package inspector

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bisohns/saido/driver"
	log "github.com/sirupsen/logrus"
)

// TcpMetrics : Metrics obtained by tcp monitoring on darwin
type TcpMetrics struct {
	// Ports map a port to a status string
	// e.g {8081: "LISTEN"}
	Ports map[int]string
}

type TcpDarwin struct {
	Command string
	Driver  *driver.Driver
	Values  TcpMetrics
}

type TcpLinux struct {
	Command string
	Driver  *driver.Driver
	Values  TcpMetrics
}

/* Parse : parsing the following kind of output
Active Internet connections (including servers)
Proto Recv-Q Send-Q  Local Address          Foreign Address        (state)
tcp4       0      0  127.0.0.1.53300        127.0.0.1.59972        ESTABLISHED
tcp4       0      0  192.168.1.172.59964    162.247.243.147.443    SYN_SENT
tcp4       0      0  192.168.1.172.59931    13.224.227.146.443     ESTABLISHED
tcp4       0      0  127.0.0.1.59905        127.0.0.1.53300        CLOSE_WAIT
*/
func (i *TcpDarwin) Parse(output string) {
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
				log.Fatal("Could not parse port number in TcpDarwin")
			}
			ports[port] = status

		}
	}
	i.Values.Ports = ports
}

func (i *TcpDarwin) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsDarwin {
		panic("Cannot use TcpDarwin on drivers outside (darwin)")
	}
	i.Driver = driver
}

func (i TcpDarwin) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *TcpDarwin) Execute() {
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
func (i *TcpLinux) Parse(output string) {
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
				log.Fatal("Could not parse port number in TcpLinux")
			}
			ports[port] = status

		}
	}
	i.Values.Ports = ports
}

func (i *TcpLinux) SetDriver(driver *driver.Driver) {
	details := (*driver).GetDetails()
	if !details.IsLinux {
		panic("Cannot use TcpLinux on drivers outside (linux)")
	}
	i.Driver = driver
}

func (i TcpLinux) driverExec() driver.Command {
	return (*i.Driver).RunCommand
}

func (i *TcpLinux) Execute() {
	output, err := i.driverExec()(i.Command)
	if err == nil {
		i.Parse(output)
	}
}

// NewTcp: Initialize a new Tcp instance
func NewTcp(driver *driver.Driver, _ ...string) (Inspector, error) {
	var tcp Inspector
	details := (*driver).GetDetails()
	if !(details.IsLinux || details.IsDarwin || details.IsWindows) {
		return nil, errors.New("Cannot use Tcp on drivers outside (linux, darwin, windows)")
	}
	if details.IsDarwin {
		tcp = &TcpDarwin{
			Command: `netstat -anp tcp`,
		}
	} else if details.IsLinux {
		tcp = &TcpLinux{
			Command: `ss -tan`,
		}
	}
	tcp.SetDriver(driver)
	return tcp, nil
}
