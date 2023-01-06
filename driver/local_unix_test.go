//go:build !windows
// +build !windows

package driver

import (
	"strings"
	"testing"
)

func TestUnixLocalRunCommand(t *testing.T) {
	d := Local{}
	output, err := d.RunCommand(`ps -A`)
	if err != nil || !strings.Contains(output, "PID") {
		t.Error(err)
	}
}

func TestUnixLocalSystemDetails(t *testing.T) {
	d := Local{}
	details, err := d.GetDetails()
	if err != nil || !(details.IsLinux || details.IsDarwin) {
		t.Errorf("Expected Darwin or Linux on unix test, got %s", details.Name)
	}
}
