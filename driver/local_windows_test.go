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

func TestWindowsLocalSystemDetails(t *testing.T) {
	d := Local{}
	details, err := d.GetDetails()
	if err != nil || !details.IsWindows {
		t.Errorf("Expected windows got %s", details.Name)
	}
}
