package driver

import (
	"strings"
	"testing"
)

func NewSSHForTest() Driver {
	return &SSH{
		User:            "dev",
		Host:            "127.0.0.1",
		Port:            2222,
		KeyFile:         "/home/diretnan/.ssh/id_rsa",
		KeyPass:         "",
		CheckKnownHosts: false,
		driverBase: driverBase{
			PollInterval: 5,
		},
	}
}

func TestSSHRunCommand(t *testing.T) {
	d := NewSSHForTest()
	output, err := d.RunCommand(`ps -A`)
	if err != nil || !strings.Contains(output, "PID") {
		t.Error(err)
	}
}

func TestSSHSystemDetails(t *testing.T) {
	d := NewSSHForTest()
	details := d.GetDetails()
	if !details.IsLinux {
		t.Errorf("Expected linux server for ssh test got %s", details.Name)
	}
}
