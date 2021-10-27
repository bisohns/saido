// +build !windows

package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bisoncorps/saido/driver"
	"github.com/bisoncorps/saido/inspector"
)

func TestDFonLocal(t *testing.T) {
	d := driver.Local{}
	// can either use NewDF() or get the interface and perform type assertion
	i := (inspector.GetInspector(`disk`)).(*inspector.DF)
	output, err := d.RunCommand(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
	if i.Values[0].Used == 0 {
		t.Error("showing percent used as 0")
	}
	fmt.Printf(`%#v`, i.Values)
}

func TestMemInfoonLocal(t *testing.T) {
	d := driver.Local{}
	// can either use NewDF() or get the interface and perform type assertion
	i := (inspector.GetInspector(`meminfo`)).(*inspector.MemInfo)
	output, err := d.ReadFile(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
	if i.Values.MemTotal == 0 {
		t.Error("showing percent used as 0")
	}
	fmt.Printf(`%#v`, i.Values)
}

func TestDockerStatsonLocal(t *testing.T) {
	d := driver.Local{}
	// can either use NewDF() or get the interface and perform type assertion
	i := (inspector.GetInspector(`dockerstats`)).(*inspector.DockerStats)
	output, err := d.RunCommand(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
	if len(i.Values) == 0 {
		t.Error("showing no running container")
	}
	fmt.Printf(`%#v`, i.Values)
}

func TestProcessonLocal(t *testing.T) {
	d := driver.Local{}
	i := inspector.NewProcess()
	output, err := d.RunCommand(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
	if len(i.Values) <= 2 {
		t.Error(err)
	}
	// Track just root PID of 1
	i.TrackPID = 1
	i.Parse(output)
	if len(i.Values) != 1 {
		t.Error("unexpected size of single PID tracking")
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
	d := driver.Local{}
	i := inspector.NewCustom(`echo /test/test`)
	output, err := d.RunCommand(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
	if strings.TrimSpace(i.Values.Output) != "/test/test" {
		t.Errorf("%s", i.Values.Output)
	}
}

func TestLoadAvgonLocal(t *testing.T) {
	d := driver.NewSSHForTest()
	i := inspector.NewLoadAvg()
	output, err := d.ReadFile(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
	if i.Values.Load1M == 0 {
		t.Errorf("%f", i.Values.Load1M)
	}
	fmt.Printf("%#v", i.Values)
}
