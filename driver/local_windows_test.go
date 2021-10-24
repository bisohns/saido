package driver

import (
	"strings"
	"testing"
)

func TestWindowsRunCommand(t *testing.T) {
	d := Local{}
	output, err := d.RunCommand(`tasklist`)
	if err != nil || !strings.Contains(output, "PID") {
		t.Error(err)
	}
}

func TestWindowsLocalGetDetails(t *testing.T) {
	d := Local{}
	output := d.GetDetails()
	if output != "Local - windows" {
		t.Error(output)
	}
}
