// +build !windows

package inspector

import (
	"testing"
)

func TestMemInfo(t *testing.T) {
	d := NewMemInfo()
	if d.Type != File || d.FilePath != `/proc/meminfo` {
		t.Error("Initialized meminfo wrongly")
	}
}
