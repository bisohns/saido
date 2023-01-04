package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bisohns/saido/driver"
	"github.com/bisohns/saido/inspector"
)

func NewLocalForTest() driver.Driver {
	return &driver.Local{}
}

func TestProcessonLocal(t *testing.T) {
	d := NewLocalForTest()
	i, _ := inspector.NewProcess(&d)
	i.Execute()
	iConcreteWin, ok := i.(*inspector.ProcessWin)
	if ok {
		if len(iConcreteWin.Values) <= 2 {
			t.Error("Less than two processes running")
		}
		if process := iConcreteWin.Values[0].Command; process != "System Idle Process" {
			t.Errorf("Expected System Idle Process as first process, found %s", iConcreteWin.Values[0].Command)
		}
	}
}

func TestCustomonLocal(t *testing.T) {
	d := NewLocalForTest()
	dfConcrete, _ := d.(*driver.Local)
	dfConcrete.EnvVars = []string{"EXAMPLES=true"}
	d = dfConcrete
	i, _ := inspector.Init(`custom`, &d, `echo %EXAMPLES%`)
	i.Execute()
	iConcrete, ok := i.(*inspector.Custom)
	if ok {
		if strings.TrimSpace(iConcrete.Values.Output) != "true" {
			t.Errorf("Expected 'true', found %s", iConcrete.Values.Output)
		}
	}
}

func TestUptimeonLocal(t *testing.T) {
	d := NewLocalForTest()
	i, _ := inspector.Init(`uptime`, &d)
	i.Execute()
	iConcrete, ok := i.(*inspector.UptimeWindows)
	if ok {
		if iConcrete.Values.Up == 0 {
			t.Error("Expected uptime on windows to be > 0")
		}
	}
}

func TestLoadAverageonLocal(t *testing.T) {
	d := NewLocalForTest()
	i, _ := inspector.Init(`loadavg`, &d)
	i.Execute()
	iConcrete, ok := i.(*inspector.LoadAvgWin)
	if ok {
		if iConcrete.Values.Load1M == 0 && !SkipNonLinuxOnCI() {
			t.Error("Expected load on windows to be > 0")
		}
	}
}

func TestMemInfoonLocal(t *testing.T) {
	d := NewLocalForTest()
	i, _ := inspector.Init(`memory`, &d)
	i.Execute()
	iConcrete, ok := i.(*inspector.MemInfoWin)
	if ok {
		fmt.Printf("%#v", iConcrete.Values)
		if iConcrete.Values.MemTotal == 0 {
			t.Error("RAM reported as empty")
		}
	}
}

func TestDFonLocal(t *testing.T) {
	d := NewLocalForTest()
	i, _ := inspector.Init(`disk`, &d)
	i.Execute()
	iConcrete, ok := i.(*inspector.DFWin)
	if ok {
		fmt.Printf("%#v", iConcrete.Values)
		if len(iConcrete.Values) < 1 {
			t.Error("DFWin not showing at least one drive")
		}
	}
}

func TestTemponLocal(t *testing.T) {
	d := NewLocalForTest()
	i, _ := inspector.Init(`temp`, &d)
	i.Execute()
	iConcrete, ok := i.(*inspector.TempWin)
	if ok {
		fmt.Printf("%#v", iConcrete.Values)
		if iConcrete.Values == nil {
			t.Error("TempWin not set on Windows")
		}
	}

}

func TestTcponLocal(t *testing.T) {
	d := NewLocalForTest()
	i, _ := inspector.Init(`tcp`, &d)
	i.Execute()
	iConcreteWindows, ok := i.(*inspector.TcpWin)
	if ok {
		if len(iConcreteWindows.Values.Ports) == 0 {
			t.Errorf("%#v", iConcreteWindows.Values.Ports)
		}
		fmt.Printf("%#v", iConcreteWindows.Values.Ports)
	}
}
