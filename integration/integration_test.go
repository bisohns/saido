package integration

import (
	"fmt"
	"testing"

	"github.com/bisoncorps/saido/driver"
	"github.com/bisoncorps/saido/inspector"
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
