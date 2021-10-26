package driver

import (
	"strings"
	"testing"
)

func TestSSHRunCommand(t *testing.T) {
	d := NewSSHForTest()
	output, err := d.RunCommand(`ps -A`)
	if err != nil || !strings.Contains(output, "PID") {
		t.Error(err)
	}
}
