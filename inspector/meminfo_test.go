// +build !windows

package inspector

import (
	"testing"
)

func TestMemInfoOnLocal(t *testing.T) {
	driver := NewLocalForTest()
	d := NewMemInfo(&driver)
	d.Execute()
}
