// +build !windows

package integration

import (
	"fmt"
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
