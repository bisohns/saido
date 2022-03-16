// +build !windows

package inspector

import (
	"testing"
)

func TestMemInfoOnLocal(t *testing.T) {
	driver := NewLocalForTest()
	d := NewMemInfo(&driver)
	d.Execute()
	iConcreteLinux, ok := d.(*MemInfoLinux)
	if ok {
		if iConcreteLinux.Values == nil {
			t.Error("Values did not get set for MemInfoLinux")
		}
	}
	iConcreteDarwin, ok := d.(*MemInfoDarwin)
	if ok {
		if iConcreteDarwin.Values == nil {
			t.Error("Values did not get set for MemInfoDarwin")
		}
	}
}
