// +build !windows

package inspector

import (
	"testing"
)

func TestDF(t *testing.T) {
	d := NewDF()
	if d.Type != Command || d.Command != `df -a` {
		t.Error("Initialized df wrongly")
	}
}
