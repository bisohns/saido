//go:build !windows
// +build !windows

package inspector

import (
	"testing"
)

func TestTempOnLocal(t *testing.T) {
	driver := NewLocalForTest()
	d, _ := NewTemp(&driver)
	d.Execute()
	TempConcreteLinux, ok := d.(*TempLinux)
	if ok {
		if TempConcreteLinux.Values == nil {
			t.Error("Values did not get set for TempLinux")
		}
	}
	TempConcreteDarwin, ok := d.(*TempDarwin)
	if ok {
		if TempConcreteDarwin.Values == nil {
			t.Error("Values did not get set for TempDarwin")
		}
	}
}

func TestTempOnSSH(t *testing.T) {
	if SkipNonLinuxOnCI() {
		return
	}
	driver := NewSSHForTest()
	d, _ := NewTemp(&driver)
	d.Execute()
	TempConcreteLinux, ok := d.(*TempLinux)
	if ok {
		if TempConcreteLinux.Values == nil {
			t.Error("Values did not get set for TempLinux")
		}
	}
	TempConcreteDarwin, ok := d.(*TempDarwin)
	if ok {
		if TempConcreteDarwin.Values == nil {
			t.Error("Values did not get set for TempDarwin")
		}
	}
}
