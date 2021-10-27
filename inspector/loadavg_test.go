// +build !windows

package inspector

import (
	"testing"
)

func TestLoadAvg(t *testing.T) {
	d := NewLoadAvg()
	if d.Type != File || d.FilePath != `/proc/loadavg` {
		t.Error("Initialized loadavg wrongly")
	}
}
