package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bisohns/saido/config"
	"github.com/bisohns/saido/driver"
	"github.com/bisohns/saido/inspector"
)

func NewWebForTest() driver.Driver {
	return &driver.Web{
		URL:    "https://duckduckgo.com",
		Method: driver.GET,
	}
}

func NewSSHForTest() driver.Driver {
	workingDir, _ := os.Getwd()
	workingDir = filepath.Dir(workingDir)
	yamlPath := fmt.Sprintf("%s/%s", workingDir, "config-test.yaml")
	conf := config.LoadConfig(yamlPath)
	dashboardInfo := config.GetDashboardInfoConfig(conf)
	return &driver.SSH{
		User:            dashboardInfo.Hosts[0].Connection.Username,
		Host:            dashboardInfo.Hosts[0].Address,
		Port:            int(dashboardInfo.Hosts[0].Connection.Port),
		KeyFile:         dashboardInfo.Hosts[0].Connection.PrivateKeyPath,
		KeyPass:         "",
		CheckKnownHosts: false,
	}
}

func TestDFonSSH(t *testing.T) {
	d := NewSSHForTest()
	i, _ := inspector.Init(`disk`, &d)
	i.Execute()
	iConcrete, ok := i.(*inspector.DF)
	if ok {
		fmt.Printf(`%#v`, iConcrete.Values)
	}
}

func TestMemInfoonSSH(t *testing.T) {
	d := NewSSHForTest()
	i, _ := inspector.NewMemInfo(&d)
	i.Execute()
	iConcreteLinux, ok := i.(*inspector.MemInfoLinux)
	if ok {
		if iConcreteLinux.Values.MemTotal == 0 {
			t.Error("showing percent used as 0")
		}
		fmt.Printf(`%#v`, iConcreteLinux.Values)
	}
	iConcreteDarwin, ok := i.(*inspector.MemInfoDarwin)
	if ok {
		if iConcreteDarwin.Values.MemTotal == 0 {
			t.Error("showing percent used as 0")
		}
		fmt.Printf(`%#v`, iConcreteDarwin.Values)
	}
}

func TestResponseTimeonWeb(t *testing.T) {
	d := NewWebForTest()
	i, _ := inspector.NewResponseTime(&d)
	i.Execute()
	iConcrete, ok := i.(*inspector.ResponseTime)
	if ok {
		if iConcrete.Values.Seconds == 0 {
			t.Error("showing response time as 0")
		}
		fmt.Printf(`%#v`, iConcrete.Values)
	}
}

func TestProcessonSSH(t *testing.T) {
	d := NewSSHForTest()
	i, _ := inspector.NewProcess(&d)
	i.Execute()
	iConcreteUnix, ok := i.(*inspector.Process)
	if ok {
		if len(iConcreteUnix.Values) <= 2 {
			t.Error("Less than two processes running")
		}
	}
	iConcreteWin, ok := i.(*inspector.ProcessWin)
	if ok {
		if len(iConcreteWin.Values) <= 2 {
			t.Error("Less than two processes running")
		}
	}
}

func TestCustomonSSH(t *testing.T) {
	d := NewSSHForTest()
	// set vars
	dfConcrete, _ := d.(*driver.SSH)
	dfConcrete.EnvVars = []string{"MONKEY=true"}
	d = dfConcrete
	i, _ := inspector.NewCustom(&d, `echo $MONKEY`)
	i.Execute()
	iConcrete, ok := i.(*inspector.Custom)
	if ok {
		if strings.TrimSpace(iConcrete.Values.Output) != "true" {
			t.Errorf("%s", iConcrete.Values.Output)
		}
	}
}

func TestLoadAvgonSSH(t *testing.T) {
	d := NewSSHForTest()
	i, _ := inspector.NewLoadAvg(&d)
	i.Execute()
	iConcreteLinux, ok := i.(*inspector.LoadAvgLinux)
	if ok {
		if iConcreteLinux.Values.Load1M == 0 {
			t.Errorf("%f", iConcreteLinux.Values.Load1M)
		}
	}
	iConcreteDarwin, ok := i.(*inspector.LoadAvgDarwin)
	if ok {
		if iConcreteDarwin.Values.Load1M == 0 {
			t.Errorf("%f", iConcreteDarwin.Values.Load1M)
		}
	}
}

func TestCustomonWeb(t *testing.T) {
	d := NewWebForTest()
	_, err := inspector.Init(`custom`, &d, `custom-command`)
	if err == nil {
		t.Error("should not instantiate custom on web")
	}
}

func TestUptimeonSSH(t *testing.T) {
	d := NewSSHForTest()
	i, _ := inspector.NewUptime(&d)
	i.Execute()
	iConcreteLinux, ok := i.(*inspector.UptimeLinux)
	if ok {
		if iConcreteLinux.Values.Up == 0 {
			t.Errorf("%f", iConcreteLinux.Values.Up)
		}
	}
	iConcreteDarwin, ok := i.(*inspector.UptimeDarwin)
	if ok {
		if iConcreteDarwin.Values.Up == 0 {
			t.Errorf("%f", iConcreteDarwin.Values.Up)
		}
	}
}
