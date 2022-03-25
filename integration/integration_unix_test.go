//go:build !windows
// +build !windows

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

func TestDFonLocal(t *testing.T) {
	d := NewLocalForTest()
	// can either use NewDF() or get the interface and perform type assertion
	i, _ := inspector.NewDF(&d)
	i.Execute()
	iConcrete, _ := i.(*inspector.DF)
	if iConcrete.Values[0].Used == 0 {
		t.Error("showing percent used as 0")
	}
	fmt.Printf(`%#v`, iConcrete.Values)
}

func TestMemInfoonLocal(t *testing.T) {
	d := NewLocalForTest()
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

func TestDockerStatsonLocal(t *testing.T) {
	if SkipNonLinuxOnCI() {
		return
	}
	d := NewLocalForTest()
	i, _ := inspector.NewDockerStats(&d)
	i.Execute()
	iConcrete, _ := i.(*inspector.DockerStats)
	if len(iConcrete.Values) == 0 {
		t.Error("showing no running container")
	}
	fmt.Printf(`%#v`, iConcrete.Values)
}

func TestProcessonLocal(t *testing.T) {
	d := NewLocalForTest()
	i, _ := inspector.NewProcess(&d)
	i.Execute()
	iConcreteUnix, ok := i.(*inspector.Process)
	if ok {
		if len(iConcreteUnix.Values) <= 2 {
			t.Error("Values are less than or equal 2")
		}
		// Track just root PID of 1
		iConcreteUnix.TrackPID = 1
		iConcreteUnix.Execute()
		if len(iConcreteUnix.Values) != 1 {
			t.Error("unexpected size of single PID tracking")
		}
	}
	iConcreteWin, ok := i.(*inspector.ProcessWin)
	if ok {
		if len(iConcreteWin.Values) <= 2 {
			t.Error("Values are less than or equal 2")
		}
		// Track just system PID of 4
		iConcreteWin.TrackPID = 4
		iConcreteWin.Execute()
		if len(iConcreteWin.Values) != 1 {
			t.Error("unexpected size of single PID tracking")
		}
	}
}

//FIXME: faulty shell globbing using custom commands
//func TestCustomonLocal(t *testing.T) {
//  d := driver.Local{
//    Vars: []string{"MONKEY=true"},
//  }
//  i := inspector.NewCustom(`echo $MONKEY`)
//  output, err := d.RunCommand(i.String())
//  if err != nil {
//    t.Error(err)
//  }
//  i.Parse(output)
//  if strings.TrimSpace(i.Values.Output) != "true" {
//    t.Errorf("%s", i.Values.Output)
//  }
//}

func TestCustomonLocal(t *testing.T) {
	d := NewLocalForTest()
	i, _ := inspector.NewCustom(&d, `echo /test/test`)
	i.Execute()
	iConcrete, _ := i.(*inspector.Custom)
	if strings.TrimSpace(iConcrete.Values.Output) != "/test/test" {
		t.Errorf("%s", iConcrete.Values.Output)
	}
}

func TestLoadAvgonLocal(t *testing.T) {
	d := NewLocalForTest()
	i, _ := inspector.NewLoadAvg(&d)
	i.Execute()
	iConcreteDarwin, ok := i.(*inspector.LoadAvgDarwin)
	if ok {
		if iConcreteDarwin.Values.Load1M == 0 {
			t.Errorf("%f", iConcreteDarwin.Values.Load1M)
		}
		fmt.Printf("%#v", iConcreteDarwin.Values)
	}
	iConcrete, ok := i.(*inspector.LoadAvgLinux)
	if ok {
		if iConcrete.Values.Load1M == 0 {
			t.Errorf("%f", iConcrete.Values.Load1M)
		}
		fmt.Printf("%#v", iConcrete.Values)
	}
}

func TestUptimeonLocal(t *testing.T) {
	d := NewLocalForTest()
	i, _ := inspector.NewUptime(&d)
	i.Execute()
	iConcreteLinux, ok := i.(*inspector.UptimeLinux)
	if ok {
		if iConcreteLinux.Values.Up == 0 {
			t.Errorf("%f", iConcreteLinux.Values.Up)
		}
		fmt.Printf("%#v", iConcreteLinux.Values)
	}
	iConcreteDarwin, ok := i.(*inspector.UptimeDarwin)
	if ok {
		if iConcreteDarwin.Values.Up == 0 {
			t.Errorf("%f", iConcreteDarwin.Values.Up)
		}
		fmt.Printf("%#v", iConcreteDarwin.Values)
	}
}

func TestTcponLocal(t *testing.T) {
	d := NewLocalForTest()
	i, _ := inspector.Init(`tcp`, &d)
	i.Execute()
	iConcreteDarwin, ok := i.(*inspector.TcpDarwin)
	if ok {
		if len(iConcreteDarwin.Values.Ports) == 0 {
			t.Errorf("%#v", iConcreteDarwin.Values.Ports)
		}
		fmt.Printf("%#v", iConcreteDarwin.Values.Ports)
	}
	iConcreteLinux, ok := i.(*inspector.TcpLinux)
	if ok {
		if len(iConcreteLinux.Values.Ports) == 0 {
			t.Errorf("%#v", iConcreteLinux.Values.Ports)
		}
		fmt.Printf("%#v", iConcreteLinux.Values.Ports)
	}
}
