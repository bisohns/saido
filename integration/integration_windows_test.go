package integration

import (
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
