// +build !windows

package inspector

import (
	"testing"

	"github.com/bisohns/saido/driver"
)

func NewLocalForTest() driver.Driver {
	return &driver.Local{}
}

func TestLoadAvg(t *testing.T) {
	testDriver := NewLocalForTest()
	loadavg, _ := NewLoadAvg(&testDriver)
	loadavg.Execute()
	loadavgConcreteLinux, ok := loadavg.(*LoadAvgLinux)
	if ok {
		if loadavgConcreteLinux.Values == nil {
			t.Error("Load metrics for linux did not get set")
		}
	}
	loadavgConcreteDarwin, ok := loadavg.(*LoadAvgDarwin)
	if ok {
		if loadavgConcreteDarwin.Values == nil {
			t.Error("Load metrics for darwin did not get set")
		}
	}
}
