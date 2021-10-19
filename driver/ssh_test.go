package driver

import (
	"strings"
	"testing"
)

func TestSSHRunCommand(t *testing.T) {
	d := SSH{
		User:            "dev",
		Host:            "127.0.0.1",
		Port:            2222,
		PubKeyFile:      "/home/runner/.ssh/id_rsa",
		PubKeyPass:      "",
		CheckKnownHosts: false,
	}
	output, err := d.runCommand(`ps -A`)
	if err != nil || !strings.Contains(output, "PID") {
		t.Error(err)
	}
}
