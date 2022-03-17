// +build !windows

package inspector

import (
	"testing"

	"github.com/bisohns/saido/driver"
)

func NewSSHForTest() driver.Driver {
	return &driver.SSH{
		User:            "dev",
		Host:            "127.0.0.1",
		Port:            2222,
		KeyFile:         "/home/diretnan/.ssh/id_rsa",
		KeyPass:         "",
		CheckKnownHosts: false,
	}
}

func TestDFOnLocal(t *testing.T) {
	driver := NewLocalForTest()
	d, _ := NewDF(&driver)
	d.Execute()
	dfConcrete, _ := d.(*DF)
	if len(dfConcrete.Values) == 0 {
		t.Error("Values are empty!")
	}
}

func TestDFOnSSH(t *testing.T) {
	driver := NewSSHForTest()
	d, _ := NewDF(&driver)
	d.Execute()
	dfConcrete, _ := d.(*DF)
	if len(dfConcrete.Values) == 0 {
		t.Error("Values are empty!")
	}
}
