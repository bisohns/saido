// +build !windows

package driver

import (
	"github.com/bisoncorps/saido/inspector"
	"strings"
	"testing"
)

func TestUnixLocalRunCommand(t *testing.T) {
	d := Local{}
	d.Supported = []inspector.Inspector{}
	output, err := d.runCommand(`ps -A`)
	if err != nil || !strings.Contains(output, "PID") {
		t.Error(err)
	}
}
