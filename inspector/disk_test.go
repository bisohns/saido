//go:build !windows
// +build !windows

package inspector

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bisohns/saido/config"
	"github.com/bisohns/saido/driver"
)

func SkipNonLinuxOnCI() bool {
	if os.Getenv("CI") == "true" {
		if runtime.GOOS != "linux" {
			return true
		}
	}
	return false
}

func NewSSHForTest() driver.Driver {
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

	return &driver.SSH{
		User:            host.Connection.Username,
		Host:            host.Address,
		Port:            int(host.Connection.Port),
		KeyFile:         host.Connection.PrivateKeyPath,
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
	if SkipNonLinuxOnCI() {
		return
	}
	driver := NewSSHForTest()
	d, _ := NewDF(&driver)
	d.Execute()
	dfConcrete, _ := d.(*DF)
	if len(dfConcrete.Values) == 0 {
		t.Error("Values are empty!")
	}
}
