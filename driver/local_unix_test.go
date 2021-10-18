// +build !windows

package driver

import (
	"strings"
	"testing"
)

func TestUnixLocalRunCommand(t *testing.T) {
	d := Local{}
	output, err := d.runCommand(`ps -A`)
	if err != nil || !strings.Contains(output, "PID") {
		t.Error(err)
	}
}
