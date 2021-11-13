package integration

import (
	"fmt"
	"testing"

	"github.com/bisohns/saido/driver"
	"github.com/bisohns/saido/inspector"
)

func TestTasklistonLocal(t *testing.T) {
	d := driver.Local{}
	i := inspector.NewTasklist()
	output, err := d.RunCommand(i.String())
	if err != nil {
		t.Error(err)
	}
	i.Parse(output)
	if len(i.Values) <= 1 {
		t.Error("showing 1 or less tasks/processes")
	}
	fmt.Printf(`%#v`, i.Values)
}
