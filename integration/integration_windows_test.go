package integration

import (
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
	dfConcrete.Vars = []string{"EXAMPLES=true"}
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
