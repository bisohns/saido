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
		PubKeyFile:      "/home/deven/.ssh/id_rsa",
		PubKeyPass:      "",
		CheckKnownHosts: false,
		fields: fields{
			PollInterval: 5,
		},
	}
	output, err := d.RunCommand(`ps -A`)
	if err != nil || !strings.Contains(output, "PID") {
		t.Error(err)
	}
}
