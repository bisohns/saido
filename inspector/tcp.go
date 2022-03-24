// Check if TCP ports are open or not
package inspector

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bisohns/saido/driver"
	log "github.com/sirupsen/logrus"
)

// TcpMetricsDarwin : Metrics obtained by tcp monitoring on darwin
type TcpMetricsDarwin struct {
	// Ports map a port to a status string
	// e.g {8081: "LISTEN"}
	Ports map[int]string
}

type TcpDarwin struct {
	Command string
	Driver  *driver.Driver
	Values  TcpMetricsDarwin
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
	}
	tcp.SetDriver(driver)
	return tcp, nil
}
