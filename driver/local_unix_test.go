// +build !windows

package driver

import (
	"strings"
	"testing"
)

func TestUnixLocalRunCommand(t *testing.T) {
	d := Local{}
	d.Supported = []string{}
	output, err := d.RunCommand(`ps -A`)
	if err != nil || !strings.Contains(output, "PID") {
		t.Error(err)
	}
}
