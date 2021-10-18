package driver

import "testing"

func TestSSHRunCommand(t *testing.T) {
	d := SSH{
		User:            "root",
		Host:            "172.17.0.2",
		Port:            2222,
		PubKeyFile:      "/home/deven/.ssh/id_rsa",
		PubKeyPass:      "",
		CheckKnownHosts: false,
	}
	_, err := d.runCommand(`ps -A`)
	if err != nil {
		t.Error(err)
	}
}
