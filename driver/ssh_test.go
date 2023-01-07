package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/bisohns/saido/config"
)

func SkipNonLinuxOnCI() bool {
	if os.Getenv("CI") == "true" {
		if runtime.GOOS != "linux" {
			return true
		}
	}
	return false
}

func NewSSHForTest() Driver {
	workingDir, _ := os.Getwd()
	workingDir = filepath.Dir(workingDir)
	yamlPath := fmt.Sprintf("%s/%s", workingDir, "config-test.yaml")
	conf := config.LoadConfig(yamlPath)
	dashboardInfo := config.GetDashboardInfoConfig(conf)
	var host config.Host
	for ind := range dashboardInfo.Hosts {
		if dashboardInfo.Hosts[ind].Address == "0.0.0.0" {
			host = dashboardInfo.Hosts[ind]
		}
	}
	return &SSH{
		User:            host.Connection.Username,
		Host:            host.Address,
		Port:            int(host.Connection.Port),
		KeyFile:         host.Connection.PrivateKeyPath,
		KeyPass:         "",
		CheckKnownHosts: false,
	}
}

func TestSSHRunCommand(t *testing.T) {
	if SkipNonLinuxOnCI() {
		return
	}
	d := NewSSHForTest()
	output, err := d.RunCommand(`ps -A`)
	if err != nil || !strings.Contains(output, "PID") {
		t.Error(err)
	}
}

func TestSSHSystemDetails(t *testing.T) {
	if SkipNonLinuxOnCI() {
		return
	}
	d := NewSSHForTest()
	details, err := d.GetDetails()
	if err != nil || !details.IsLinux {
		t.Errorf("Expected linux server for ssh test got %#v", details)
	}
}
