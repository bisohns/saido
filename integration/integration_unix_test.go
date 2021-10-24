// +build !windows

package integration

import (
	"github.com/bisoncorps/saido/driver"
	"github.com/bisoncorps/saido/inspector"
	"testing"
)

func TestDFonLocal(t *testing.T) {
	d := driver.Local{}
	i := inspector.NewDF()
	output, err := d.RunCommand(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
}
