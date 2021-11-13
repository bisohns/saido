package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bisohns/saido/driver"
	"github.com/bisohns/saido/inspector"
)

func TestDFonSSH(t *testing.T) {
	d := driver.NewSSHForTest()
	i := inspector.NewDF()
	output, err := d.RunCommand(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
	fmt.Printf(`%#v`, i.Values)
}

func TestMemInfoonSSH(t *testing.T) {
	d := driver.NewSSHForTest()
	i := inspector.NewMemInfo()
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

func TestResponseTimeonWeb(t *testing.T) {
	d := driver.NewWebForTest()
	i := inspector.NewResponseTime()
	output, err := d.RunCommand(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
	if i.Values.Seconds == 0 {
		t.Error("showing response time as 0")
	}
	fmt.Printf(`%#v`, i.Values)
}

func TestProcessonSSH(t *testing.T) {
	d := driver.NewSSHForTest()
	i := inspector.NewProcess()
	output, err := d.RunCommand(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
	if len(i.Values) <= 2 {
		t.Error(err)
	}
}

func TestCustomonSSH(t *testing.T) {
	d := driver.NewSSHForTest()
	// set vars
	d.Vars = []string{"MONKEY=true"}
	i := inspector.NewCustom(`echo $MONKEY`)
	output, err := d.RunCommand(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
	if strings.TrimSpace(i.Values.Output) != "true" {
		t.Errorf("%s", i.Values.Output)
	}
}

func TestLoadAvgonSSH(t *testing.T) {
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
}

func TestUptimeonSSH(t *testing.T) {
	d := driver.NewSSHForTest()
	i := inspector.NewUptime()
	output, err := d.ReadFile(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
	if i.Values.Up == 0 {
		t.Errorf("%f", i.Values.Up)
	}
}
